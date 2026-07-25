package catalog

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"platform/gokit/errs"
	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/model"
	"platform/products/nav/api/internal/naverr"
	"platform/products/nav/api/internal/navprofile"
)

type countingChecker struct {
	mu      sync.Mutex
	current int
	maximum int
}

func (c *countingChecker) Check(context.Context, string) model.LinkHealth {
	c.mu.Lock()
	c.current++
	c.maximum = max(c.maximum, c.current)
	c.mu.Unlock()
	time.Sleep(5 * time.Millisecond)
	c.mu.Lock()
	c.current--
	c.mu.Unlock()
	return model.LinkHealth{Status: "healthy", HTTPStatus: 200}
}

type fakeStore struct {
	links          []*model.Link
	lastLinkFilter dao.LinkFilter
	existing       map[string]bool
	settings       *model.SiteSettings
	health         map[string]model.LinkHealth
	healthMu       sync.Mutex
	clicks         map[string]int
}

func (f *fakeStore) Categories(context.Context) ([]*model.Category, error) {
	return []*model.Category{{ID: "develop", Title: "开发工程"}}, nil
}

func (f *fakeStore) Groups(context.Context) ([]*model.Group, error) {
	return []*model.Group{{ID: "references", CategoryID: "develop", Title: "文档与规范"}}, nil
}

func (f *fakeStore) Links(_ context.Context, filter dao.LinkFilter) ([]*model.Link, error) {
	f.lastLinkFilter = filter
	return f.links, nil
}

func (f *fakeStore) LinksByIDs(_ context.Context, ids []string) ([]*model.Link, error) {
	links := make([]*model.Link, 0, len(ids))
	for _, id := range ids {
		for _, link := range f.links {
			if link.ID == id {
				links = append(links, link)
			}
		}
	}
	return links, nil
}

func (f *fakeStore) LinkByID(_ context.Context, id string) (*model.Link, error) {
	for _, link := range f.links {
		if link.ID == id {
			return link, nil
		}
	}
	return nil, nil
}

func (f *fakeStore) CountLinks(context.Context, dao.LinkFilter) (int, error) {
	return len(f.links), nil
}
func (f *fakeStore) LinkStatusCounts(context.Context, dao.LinkFilter) (map[string]int, error) {
	return map[string]int{"all": len(f.links)}, nil
}
func (f *fakeStore) LinkHealthCounts(context.Context) (map[string]int, error) {
	return map[string]int{"all": len(f.links)}, nil
}
func (f *fakeStore) Tags(context.Context, string) ([]*model.Tag, error) { return nil, nil }

func (f *fakeStore) GroupBelongsToCategory(_ context.Context, groupID, categoryID string) (bool, error) {
	return groupID == "references" && categoryID == "develop", nil
}

func (f *fakeStore) LinkExists(_ context.Context, id string) (bool, error) {
	return f.existing[id], nil
}

func (f *fakeStore) RecordClick(_ context.Context, id string) (bool, error) {
	if !f.existing[id] {
		return false, nil
	}
	if f.clicks == nil {
		f.clicks = map[string]int{}
	}
	f.clicks[id]++
	return true, nil
}

func (f *fakeStore) UpdateLinkHealth(_ context.Context, id string, health model.LinkHealth) (bool, error) {
	f.healthMu.Lock()
	defer f.healthMu.Unlock()
	if f.health == nil {
		f.health = map[string]model.LinkHealth{}
	}
	f.health[id] = health
	return true, nil
}

func (f *fakeStore) InsertLink(_ context.Context, link *model.Link) error {
	f.links = append(f.links, link)
	return nil
}

func (f *fakeStore) UpdateLink(context.Context, *model.Link) (bool, error) {
	return true, nil
}

func (f *fakeStore) DeleteLink(context.Context, string) (bool, error) {
	return true, nil
}

func (f *fakeStore) BulkUpdateLinks(context.Context, []string, string) (int, error) { return 1, nil }
func (f *fakeStore) BulkDeleteLinks(context.Context, []string) (int, error)         { return 1, nil }
func (f *fakeStore) InsertCategory(context.Context, *model.Category) error          { return nil }
func (f *fakeStore) UpdateCategory(context.Context, *model.Category) (bool, error)  { return true, nil }
func (f *fakeStore) DeleteCategory(context.Context, string) (bool, error)           { return true, nil }
func (f *fakeStore) InsertGroup(context.Context, *model.Group) error                { return nil }
func (f *fakeStore) UpdateGroup(context.Context, *model.Group) (bool, error)        { return true, nil }
func (f *fakeStore) DeleteGroup(context.Context, string) (bool, error)              { return true, nil }
func (f *fakeStore) RenameTag(context.Context, string, string) (int, error)         { return 1, nil }
func (f *fakeStore) DeleteTag(context.Context, string) (int, error)                 { return 1, nil }
func (f *fakeStore) SiteSettings(context.Context) (*model.SiteSettings, error) {
	return f.settings, nil
}
func (f *fakeStore) UpsertSiteSettings(context.Context, *model.SiteSettings) error { return nil }

func TestSettingsRequiresProvisionedConfiguration(t *testing.T) {
	service := New(&fakeStore{}, Site{Name: "compiled fallback must not be used"})
	service.SetSiteProfile(navprofile.NewMemory())
	_, err := service.PublicSite(context.Background())
	var coded *errs.Coded
	if !errors.As(err, &coded) || coded.Code != naverr.CodeNotInitialized {
		t.Fatalf("error = %#v, want nav.not_initialized", err)
	}
}

func TestBulkLinksReportsMissingIDs(t *testing.T) {
	store := &fakeStore{existing: map[string]bool{"present": true}}
	result, err := New(store, Site{}).BulkLinks(context.Background(), []string{"present", "missing"}, "publish")
	if err != nil {
		t.Fatal(err)
	}
	if result.Changed != 1 || len(result.FailedIDs) != 1 || result.FailedIDs[0] != "missing" {
		t.Fatalf("unexpected bulk result: %#v", result)
	}
}

func TestRunChecksCapsConcurrency(t *testing.T) {
	store := &fakeStore{}
	ids := make([]string, 20)
	for index := range ids {
		ids[index] = fmt.Sprintf("link-%02d", index)
		store.links = append(store.links, &model.Link{ID: ids[index], URL: "https://example.com"})
	}
	checker := &countingChecker{}
	service := New(store, Site{})
	service.checker = checker
	results, err := service.RunChecks(context.Background(), ids)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(ids) || checker.maximum < 2 || checker.maximum > 12 {
		t.Fatalf("results=%d maximum=%d, want 20 results and concurrency 2..12", len(results), checker.maximum)
	}
}

func TestRunChecksRejectsMoreThanFiftySelectedIDs(t *testing.T) {
	ids := make([]string, 51)
	for index := range ids {
		ids[index] = fmt.Sprintf("link-%02d", index)
	}

	_, err := New(&fakeStore{}, Site{}).RunChecks(context.Background(), ids)
	if err == nil {
		t.Fatal("expected more than 50 selected ids to be rejected")
	}
}

func TestRunFilteredChecksChecksEveryMatchingLink(t *testing.T) {
	store := &fakeStore{}
	for index := range 72 {
		store.links = append(store.links, &model.Link{ID: fmt.Sprintf("link-%02d", index), URL: "https://example.com"})
	}
	service := New(store, Site{})
	service.checker = &countingChecker{}

	results, err := service.RunFilteredChecks(context.Background(), dao.LinkFilter{
		Query: "docs", Health: "unchecked", Page: 3, Size: 15,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 72 {
		t.Fatalf("results = %d, want all 72 filtered links", len(results))
	}
	if store.lastLinkFilter.Query != "docs" || store.lastLinkFilter.Health != "unchecked" || store.lastLinkFilter.Page != 0 || store.lastLinkFilter.Size != 0 {
		t.Fatalf("unexpected filter: %#v", store.lastLinkFilter)
	}
}

func TestCreateLinkNormalizesInput(t *testing.T) {
	store := &fakeStore{}
	service := New(store, Site{Name: "Nav"})
	link, err := service.CreateLink(context.Background(), LinkInput{
		CategoryID:  "develop",
		GroupID:     "references",
		Title:       "  MDN  ",
		URL:         "https://developer.mozilla.org/",
		Description: " Web docs ",
		Tags:        []string{"Web", "web", " CSS "},
		Kind:        "reference",
		Status:      "published",
	})
	if err != nil {
		t.Fatal(err)
	}
	if link.ID == "" || link.Title != "MDN" || len(link.Tags) != 2 || link.Tags[1] != "CSS" {
		t.Fatalf("unexpected normalized link: %#v", link)
	}
}

func TestCreateLinkRejectsInvalidGroup(t *testing.T) {
	service := New(&fakeStore{}, Site{Name: "Nav"})
	_, err := service.CreateLink(context.Background(), LinkInput{
		CategoryID:  "ai",
		GroupID:     "references",
		Title:       "Example",
		URL:         "https://example.com/",
		Description: "Example",
		Kind:        "official",
	})
	if err == nil {
		t.Fatal("expected invalid group error")
	}
}

func TestCreateLinkRejectsUnsafeURL(t *testing.T) {
	service := New(&fakeStore{}, Site{Name: "Nav"})
	_, err := service.CreateLink(context.Background(), LinkInput{
		CategoryID:  "develop",
		GroupID:     "references",
		Title:       "Unsafe",
		URL:         "javascript:alert(1)",
		Description: "Unsafe",
		Kind:        "tool",
	})
	if err == nil {
		t.Fatal("expected invalid URL error")
	}
	var coded *errs.Coded
	if !errors.As(err, &coded) || coded.Code != errs.CommonValidationFailed {
		t.Fatalf("error = %#v, want common.validation_failed", err)
	}
	if _, ok := coded.Params["details"]; !ok {
		t.Fatalf("validation error missing field details: %#v", coded.Params)
	}
}
