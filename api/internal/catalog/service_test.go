package catalog

import (
	"context"
	"errors"
	"testing"

	"platform/gokit/errs"
	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/model"
)

type fakeStore struct {
	links []*model.Link
}

func (f *fakeStore) Categories(context.Context) ([]*model.Category, error) {
	return []*model.Category{{ID: "develop", Title: "开发工程"}}, nil
}

func (f *fakeStore) Groups(context.Context) ([]*model.Group, error) {
	return []*model.Group{{ID: "references", CategoryID: "develop", Title: "文档与规范"}}, nil
}

func (f *fakeStore) Links(context.Context, dao.LinkFilter) ([]*model.Link, error) {
	return f.links, nil
}

func (f *fakeStore) GroupBelongsToCategory(_ context.Context, groupID, categoryID string) (bool, error) {
	return groupID == "references" && categoryID == "develop", nil
}

func (f *fakeStore) LinkExists(context.Context, string) (bool, error) {
	return false, nil
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
