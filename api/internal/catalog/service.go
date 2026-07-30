package catalog

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"github.com/yueli-official/foundation/go/audit"
	"github.com/yueli-official/foundation/go/siteprofile"

	"github.com/yueli-official/nav/api/internal/dao"
	"github.com/yueli-official/nav/api/internal/model"
	"github.com/yueli-official/nav/api/internal/navaudit"
	"github.com/yueli-official/nav/api/internal/naverr"
	"github.com/yueli-official/nav/api/internal/navprofile"
)

const StatusPublished = "published"

var allowedKinds = []string{"official", "tool", "community", "learning", "resource", "reference", "research"}
var allowedStatuses = []string{StatusPublished, "draft", "archived"}

type Store interface {
	Categories(context.Context) ([]*model.Category, error)
	Groups(context.Context) ([]*model.Group, error)
	Links(context.Context, dao.LinkFilter) ([]*model.Link, error)
	LinksByIDs(context.Context, []string) ([]*model.Link, error)
	LinkByID(context.Context, string) (*model.Link, error)
	CountLinks(context.Context, dao.LinkFilter) (int, error)
	LinkStatusCounts(context.Context, dao.LinkFilter) (map[string]int, error)
	LinkHealthCounts(context.Context) (map[string]int, error)
	Tags(context.Context, string) ([]*model.Tag, error)
	GroupBelongsToCategory(context.Context, string, string) (bool, error)
	LinkExists(context.Context, string) (bool, error)
	RecordClick(context.Context, string) (bool, error)
	UpdateLinkHealth(context.Context, string, model.LinkHealth) (bool, error)
	InsertLink(context.Context, *model.Link) error
	UpdateLink(context.Context, *model.Link) (bool, error)
	DeleteLink(context.Context, string) (bool, error)
	BulkUpdateLinks(context.Context, []string, string) (int, error)
	BulkDeleteLinks(context.Context, []string) (int, error)
	InsertCategory(context.Context, *model.Category) error
	UpdateCategory(context.Context, *model.Category) (bool, error)
	DeleteCategory(context.Context, string) (bool, error)
	InsertGroup(context.Context, *model.Group) error
	UpdateGroup(context.Context, *model.Group) (bool, error)
	DeleteGroup(context.Context, string) (bool, error)
	RenameTag(context.Context, string, string) (int, error)
	DeleteTag(context.Context, string) (int, error)
	SiteSettings(context.Context) (*model.SiteSettings, error)
	UpsertSiteSettings(context.Context, *model.SiteSettings) error
}

type Site struct {
	Revision          uint64
	RuntimeRevision   uint64
	ETag              string
	Name              string
	Title             string
	Description       string
	SearchPlaceholder string
	FooterTagline     string
}

type AdminSiteSettings struct {
	Snapshot          siteprofile.Snapshot
	Schema            siteprofile.FormSchema
	SearchPlaceholder string
	RuntimeRevision   uint64
	ETag              string
}

type profileSettingsStore interface {
	LegacySiteSettings(context.Context) (*model.LegacySiteSettings, error)
	SaveSiteSettingsWithHook(context.Context, string, uint64, dao.TransactionHook) error
	CutoverSiteSettingsWithHook(context.Context, string, dao.TransactionHook) error
}

type LinkInput struct {
	ID           string
	CategoryID   string
	GroupID      string
	Title        string
	URL          string
	Description  string
	Tags         []string
	Keywords     []string
	Kind         string
	Featured     bool
	Status       string
	SortOrder    int
	SubmitterSub string
}

type Catalog struct {
	Site       Site
	Categories []*model.Category
	Groups     []*model.Group
	Links      []*model.Link
}

type AdminLinkPage struct {
	Links  []*model.Link
	Total  int
	Counts map[string]int
	Tags   []*model.Tag
}

type GroupPage struct {
	Site     Site
	Category *model.Category
	Group    *model.Group
	Links    []*model.Link
	Total    int
	Page     int
	Size     int
}

type AdminCheckPage struct {
	Links  []*model.Link
	Total  int
	Counts map[string]int
}

type Service struct {
	store         Store
	site          Site
	checker       LinkChecker
	faviconClient *http.Client
	profiles      *navprofile.Manager
	audit         *navaudit.Journal
}

func (s *Service) SetSiteProfile(profiles *navprofile.Manager) {
	s.profiles = profiles
}

func New(store Store, site Site) *Service {
	return &Service{
		store: store, site: site, checker: NewHTTPLinkChecker(),
		faviconClient: newSafeHTTPClient(8*time.Second, true),
	}
}

func (s *Service) PublicCatalog(ctx context.Context) (*Catalog, error) {
	return s.catalog(ctx, dao.LinkFilter{Status: StatusPublished})
}

func (s *Service) PublicGroup(ctx context.Context, groupID string, page, size int, sort string) (*GroupPage, error) {
	groupID = strings.TrimSpace(groupID)
	groups, err := s.store.Groups(ctx)
	if err != nil {
		return nil, err
	}
	groupIndex := slices.IndexFunc(groups, func(group *model.Group) bool { return group.ID == groupID })
	if groupIndex < 0 {
		return nil, naverr.NotFound(groupID)
	}
	group := groups[groupIndex]
	categories, err := s.store.Categories(ctx)
	if err != nil {
		return nil, err
	}
	categoryIndex := slices.IndexFunc(categories, func(category *model.Category) bool { return category.ID == group.CategoryID })
	if categoryIndex < 0 {
		return nil, naverr.NotFound(group.CategoryID)
	}
	page, size = max(page, 1), min(max(size, 1), 60)
	filter := dao.LinkFilter{GroupID: groupID, Status: StatusPublished, Page: page, Size: size, Sort: sort}
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.store.CountLinks(ctx, filter)
	if err != nil {
		return nil, err
	}
	settings, err := s.PublicSite(ctx)
	if err != nil {
		return nil, err
	}
	return &GroupPage{
		Site:     settings,
		Category: categories[categoryIndex], Group: group, Links: links, Total: total, Page: page, Size: size,
	}, nil
}

func (s *Service) RecordClick(ctx context.Context, id string) (bool, error) {
	recorded, err := s.store.RecordClick(ctx, strings.TrimSpace(id))
	if err != nil {
		return false, err
	}
	if !recorded {
		return false, naverr.NotFound(id)
	}
	return true, nil
}

func (s *Service) AdminLinks(ctx context.Context, filter dao.LinkFilter) (*AdminLinkPage, error) {
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.store.CountLinks(ctx, filter)
	if err != nil {
		return nil, err
	}
	counts, err := s.store.LinkStatusCounts(ctx, filter)
	if err != nil {
		return nil, err
	}
	tags, err := s.store.Tags(ctx, "")
	if err != nil {
		return nil, err
	}
	return &AdminLinkPage{Links: links, Total: total, Counts: counts, Tags: tags}, nil
}

func (s *Service) AdminChecks(ctx context.Context, filter dao.LinkFilter) (*AdminCheckPage, error) {
	filter.Sort = "health"
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return nil, err
	}
	total, err := s.store.CountLinks(ctx, filter)
	if err != nil {
		return nil, err
	}
	counts, err := s.store.LinkHealthCounts(ctx)
	if err != nil {
		return nil, err
	}
	return &AdminCheckPage{Links: links, Total: total, Counts: counts}, nil
}

func (s *Service) RunChecks(ctx context.Context, ids []string) ([]*model.Link, error) {
	if len(ids) > 50 {
		return nil, naverr.Validation("ids", "maximum", map[string]any{"max": 50})
	}
	ids = normalize(ids, max(len(ids), 1))
	if len(ids) == 0 {
		return nil, naverr.Validation("ids", "required", nil)
	}
	links, err := s.store.LinksByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if len(links) != len(ids) {
		return nil, naverr.NotFound("one_or_more_links")
	}
	return s.runLinkChecks(ctx, links)
}

func (s *Service) RunFilteredChecks(ctx context.Context, filter dao.LinkFilter) ([]*model.Link, error) {
	filter.Page = 0
	filter.Size = 0
	filter.Sort = "health"
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return nil, err
	}
	return s.runLinkChecks(ctx, links)
}

func (s *Service) runLinkChecks(ctx context.Context, links []*model.Link) ([]*model.Link, error) {
	if len(links) == 0 {
		return []*model.Link{}, nil
	}
	results := make([]*model.Link, len(links))
	jobs := make(chan int)
	var (
		workers  sync.WaitGroup
		firstErr error
		errOnce  sync.Once
	)
	workerCount := min(12, len(links))
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for index := range jobs {
				link := links[index]
				health := s.checker.Check(ctx, link.URL)
				health.CheckedAt = gtime.Now()
				updated, updateErr := s.store.UpdateLinkHealth(ctx, link.ID, health)
				if updateErr != nil || !updated {
					if updateErr == nil {
						updateErr = naverr.NotFound(link.ID)
					}
					errOnce.Do(func() { firstErr = updateErr })
					continue
				}
				link.HealthStatus = health.Status
				link.HealthHTTPStatus = health.HTTPStatus
				link.HealthLatencyMS = health.LatencyMS
				link.HealthError = health.Error
				link.LastCheckedAt = health.CheckedAt
				results[index] = link
			}
		}()
	}
	for index := range links {
		jobs <- index
	}
	close(jobs)
	workers.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	return results, nil
}

func (s *Service) AdminStructure(ctx context.Context) (*Catalog, error) {
	return s.catalog(ctx, dao.LinkFilter{})
}

func (s *Service) CreateLink(ctx context.Context, input LinkInput) (*model.Link, error) {
	link, err := s.validate(ctx, input)
	if err != nil {
		return nil, err
	}
	if link.ID == "" {
		link.ID = uuid.NewString()
	}
	exists, err := s.store.LinkExists(ctx, link.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, naverr.Conflict(link.ID)
	}
	var insertErr error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		insertErr = store.InsertLinkWithHook(
			ctx, link,
			s.linkAuditHook(ctx, navigationAction(link.Status), "nav.link", link.ID, linkDigest(link), 0),
		)
	} else {
		insertErr = s.store.InsertLink(ctx, link)
	}
	if insertErr != nil {
		return nil, insertErr
	}
	created, err := s.store.LinkByID(ctx, link.ID)
	if err != nil {
		return nil, err
	}
	if created != nil {
		return created, nil
	}
	return link, nil
}

func (s *Service) Link(ctx context.Context, id string) (*model.Link, error) {
	link, err := s.store.LinkByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, naverr.NotFound(id)
	}
	return link, nil
}

func (s *Service) UpdateLink(ctx context.Context, id string, input LinkInput) (*model.Link, error) {
	input.ID = strings.TrimSpace(id)
	link, err := s.validate(ctx, input)
	if err != nil {
		return nil, err
	}
	var updated bool
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		updated, err = store.UpdateLinkWithHook(
			ctx, link,
			s.linkAuditHook(ctx, navigationAction(link.Status), "nav.link", link.ID, linkDigest(link), 0),
		)
	} else {
		updated, err = s.store.UpdateLink(ctx, link)
	}
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, naverr.NotFound(id)
	}
	current, err := s.store.LinkByID(ctx, link.ID)
	if err != nil {
		return nil, err
	}
	if current != nil {
		return current, nil
	}
	return link, nil
}

func (s *Service) DeleteLink(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	var deleted bool
	var err error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		deleted, err = store.DeleteLinkWithHook(
			ctx, id,
			s.linkAuditHook(ctx, navaudit.ActionNavigationDeleted, "nav.link", id, "", 0),
		)
	} else {
		deleted, err = s.store.DeleteLink(ctx, id)
	}
	if err != nil {
		return err
	}
	if !deleted {
		return naverr.NotFound(id)
	}
	return nil
}

type BulkResult struct {
	Changed   int
	FailedIDs []string
}

func (s *Service) BulkLinks(ctx context.Context, ids []string, action string) (BulkResult, error) {
	if len(ids) > 100 {
		return BulkResult{}, naverr.Validation("ids", "maximum", map[string]any{"max": 100})
	}
	ids = normalize(ids, max(len(ids), 1))
	if len(ids) == 0 {
		return BulkResult{}, naverr.Validation("ids", "required", nil)
	}
	statuses := map[string]string{"publish": "published", "draft": "draft", "archive": "archived"}
	status, statusAction := statuses[action]
	if action != "delete" && !statusAction {
		return BulkResult{}, naverr.Validation("action", "one_of", map[string]any{"allowed": []string{"publish", "draft", "archive", "delete"}})
	}
	eligible := make([]string, 0, len(ids))
	failed := make([]string, 0)
	for _, id := range ids {
		exists, existsErr := s.store.LinkExists(ctx, id)
		if existsErr != nil {
			return BulkResult{}, existsErr
		}
		if exists {
			eligible = append(eligible, id)
		} else {
			failed = append(failed, id)
		}
	}
	if len(eligible) == 0 {
		return BulkResult{FailedIDs: failed}, nil
	}
	var changed int
	var err error
	if action == "delete" {
		if store, ok := s.store.(auditedStore); ok && s.audit != nil {
			changed, err = store.BulkDeleteLinksWithHook(
				ctx, eligible,
				s.linkAuditHook(ctx, navaudit.ActionNavigationDeleted, "nav.link_batch", uuid.NewString(), "", len(eligible)),
			)
		} else {
			changed, err = s.store.BulkDeleteLinks(ctx, eligible)
		}
	} else {
		if store, ok := s.store.(auditedStore); ok && s.audit != nil {
			changed, err = store.BulkUpdateLinksWithHook(
				ctx, eligible, status,
				s.linkAuditHook(ctx, navigationAction(status), "nav.link_batch", uuid.NewString(), "", len(eligible)),
			)
		} else {
			changed, err = s.store.BulkUpdateLinks(ctx, eligible, status)
		}
	}
	if err != nil {
		return BulkResult{}, err
	}
	for _, id := range eligible {
		exists, existsErr := s.store.LinkExists(ctx, id)
		if existsErr != nil {
			return BulkResult{}, existsErr
		}
		if (action == "delete" && exists) || (action != "delete" && !exists) {
			failed = append(failed, id)
		}
	}
	return BulkResult{Changed: changed, FailedIDs: failed}, nil
}

func (s *Service) CreateCategory(ctx context.Context, input model.Category) (*model.Category, error) {
	category, err := validateCategory(input)
	if err != nil {
		return nil, err
	}
	category.ID = uuid.NewString()
	var insertErr error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		insertErr = store.InsertCategoryWithHook(ctx, category, s.taxonomyAuditHook(ctx, "category", category.ID, 0))
	} else {
		insertErr = s.store.InsertCategory(ctx, category)
	}
	if insertErr != nil {
		return nil, insertErr
	}
	return category, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id string, input model.Category) (*model.Category, error) {
	input.ID = strings.TrimSpace(id)
	category, err := validateCategory(input)
	if err != nil {
		return nil, err
	}
	var updated bool
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		updated, err = store.UpdateCategoryWithHook(ctx, category, s.taxonomyAuditHook(ctx, "category", category.ID, 0))
	} else {
		updated, err = s.store.UpdateCategory(ctx, category)
	}
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, naverr.NotFound(id)
	}
	return category, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	var deleted bool
	var err error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		deleted, err = store.DeleteCategoryWithHook(ctx, id, s.taxonomyAuditHook(ctx, "category", id, 0))
	} else {
		deleted, err = s.store.DeleteCategory(ctx, id)
	}
	if err != nil {
		return err
	}
	if !deleted {
		return naverr.Conflict(strings.TrimSpace(id))
	}
	return nil
}

func (s *Service) CreateGroup(ctx context.Context, input model.Group) (*model.Group, error) {
	group, err := s.validateGroup(ctx, input)
	if err != nil {
		return nil, err
	}
	group.ID = uuid.NewString()
	var insertErr error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		insertErr = store.InsertGroupWithHook(ctx, group, s.taxonomyAuditHook(ctx, "group", group.ID, 0))
	} else {
		insertErr = s.store.InsertGroup(ctx, group)
	}
	if insertErr != nil {
		return nil, insertErr
	}
	return group, nil
}

func (s *Service) UpdateGroup(ctx context.Context, id string, input model.Group) (*model.Group, error) {
	input.ID = strings.TrimSpace(id)
	group, err := s.validateGroup(ctx, input)
	if err != nil {
		return nil, err
	}
	var updated bool
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		updated, err = store.UpdateGroupWithHook(ctx, group, s.taxonomyAuditHook(ctx, "group", group.ID, 0))
	} else {
		updated, err = s.store.UpdateGroup(ctx, group)
	}
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, naverr.NotFound(id)
	}
	return group, nil
}

func (s *Service) DeleteGroup(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	var deleted bool
	var err error
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		deleted, err = store.DeleteGroupWithHook(ctx, id, s.taxonomyAuditHook(ctx, "group", id, 0))
	} else {
		deleted, err = s.store.DeleteGroup(ctx, id)
	}
	if err != nil {
		return err
	}
	if !deleted {
		return naverr.Conflict(strings.TrimSpace(id))
	}
	return nil
}

func (s *Service) Tags(ctx context.Context, query string) ([]*model.Tag, error) {
	return s.store.Tags(ctx, strings.TrimSpace(query))
}

func (s *Service) RenameTag(ctx context.Context, source, target string) (int, error) {
	source, target = strings.TrimSpace(source), strings.TrimSpace(target)
	if source == "" || target == "" {
		return 0, naverr.Validation("tag", "required", nil)
	}
	if strings.EqualFold(source, target) {
		return 0, naverr.Validation("target", "different", nil)
	}
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		return store.RenameTagWithHook(ctx, source, target, s.taxonomyAuditHook(ctx, "tag", source, 0))
	}
	return s.store.RenameTag(ctx, source, target)
}

func (s *Service) DeleteTag(ctx context.Context, name string) (int, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, naverr.Validation("name", "required", nil)
	}
	if store, ok := s.store.(auditedStore); ok && s.audit != nil {
		return store.DeleteTagWithHook(ctx, name, s.taxonomyAuditHook(ctx, "tag", name, 0))
	}
	return s.store.DeleteTag(ctx, name)
}

func (s *Service) PublicSite(ctx context.Context) (Site, error) {
	if s.profiles == nil {
		return Site{}, errors.New("navigation site profile module is not configured")
	}
	settings, err := s.store.SiteSettings(ctx)
	if err != nil {
		return Site{}, err
	}
	if settings == nil {
		return Site{}, naverr.NotInitialized("site_settings")
	}
	projection, err := s.profiles.PublicAt(ctx)
	if err != nil {
		return Site{}, err
	}
	return siteFromSnapshot(projection.Snapshot, *settings), nil
}

func (s *Service) AdminSiteSettings(ctx context.Context) (AdminSiteSettings, error) {
	if s.profiles == nil {
		return AdminSiteSettings{}, errors.New("navigation site profile module is not configured")
	}
	settings, err := s.store.SiteSettings(ctx)
	if err != nil {
		return AdminSiteSettings{}, err
	}
	if settings == nil {
		return AdminSiteSettings{}, naverr.NotInitialized("site_settings")
	}
	snapshot, err := s.profiles.Get(ctx)
	if err != nil {
		return AdminSiteSettings{}, err
	}
	runtimeRevision := normalizedRuntimeRevision(settings.RuntimeRevision)
	return AdminSiteSettings{
		Snapshot: snapshot, Schema: s.profiles.Schema(),
		SearchPlaceholder: settings.SearchPlaceholder, RuntimeRevision: runtimeRevision,
		ETag: consumerETag(snapshot.ETag, runtimeRevision, settings.SearchPlaceholder),
	}, nil
}

func (s *Service) SaveAdminSiteSettings(
	ctx context.Context,
	expected siteprofile.Revision,
	expectedRuntimeRevision uint64,
	profile siteprofile.Profile,
	searchPlaceholder string,
) (AdminSiteSettings, error) {
	searchPlaceholder = strings.TrimSpace(searchPlaceholder)
	if searchPlaceholder == "" {
		return AdminSiteSettings{}, naverr.Validation("searchPlaceholder", "required", nil)
	}
	store, ok := s.store.(profileSettingsStore)
	if !ok {
		return AdminSiteSettings{}, errors.New("navigation store does not support atomic site profile settings")
	}
	err := store.SaveSiteSettingsWithHook(ctx, searchPlaceholder, expectedRuntimeRevision, func(ctx context.Context, tx *sql.Tx) error {
		result, replaceErr := s.profiles.ReplaceTx(ctx, tx, siteprofile.ReplaceCommand{
			ExpectedRevision: expected,
			Profile:          profile,
		})
		if replaceErr != nil {
			return replaceErr
		}
		if s.audit == nil || !result.Changed {
			return nil
		}
		hook := s.audit.Hook(
			ctx, navaudit.ActionSiteProfilePublished, uuid.NewString(),
			audit.Target{Type: "nav.site_profile", ID: "default"},
			navaudit.Evidence{
				Revision: uint64(result.Snapshot.Revision),
				Digest:   string(result.Snapshot.DocumentDigest),
			},
		)
		return hook(ctx, tx)
	})
	if errors.Is(err, dao.ErrSiteSettingsRevisionConflict) {
		return AdminSiteSettings{}, naverr.RevisionConflict()
	}
	if err != nil {
		return AdminSiteSettings{}, mapSiteProfileError(err)
	}
	return s.AdminSiteSettings(ctx)
}

func (s *Service) EnsureSiteProfile(ctx context.Context) error {
	if s.profiles == nil {
		return errors.New("navigation site profile module is not configured")
	}
	if _, err := s.profiles.Get(ctx); err == nil {
		store, ok := s.store.(profileSettingsStore)
		if !ok {
			return errors.New("navigation store does not support site profile cutover")
		}
		settings, settingsErr := s.store.SiteSettings(ctx)
		if settingsErr != nil {
			return settingsErr
		}
		if settings == nil || strings.TrimSpace(settings.SearchPlaceholder) == "" {
			return naverr.NotInitialized("site_settings")
		}
		return store.CutoverSiteSettingsWithHook(ctx, settings.SearchPlaceholder, nil)
	} else if !errors.Is(err, siteprofile.ErrNotInitialized) {
		return err
	}
	store, ok := s.store.(profileSettingsStore)
	if !ok {
		return errors.New("navigation store does not support site profile cutover")
	}
	legacy, err := store.LegacySiteSettings(ctx)
	if err != nil {
		return err
	}
	if legacy == nil {
		legacy = &model.LegacySiteSettings{
			Name: s.site.Name, Title: s.site.Title, Description: s.site.Description,
			SearchPlaceholder: s.site.SearchPlaceholder, FooterTagline: s.site.FooterTagline,
		}
	}
	if strings.TrimSpace(legacy.SearchPlaceholder) == "" {
		return naverr.Validation("searchPlaceholder", "required", nil)
	}
	return store.CutoverSiteSettingsWithHook(ctx, legacy.SearchPlaceholder, func(ctx context.Context, tx *sql.Tx) error {
		_, replaceErr := s.profiles.ReplaceTx(ctx, tx, siteprofile.ReplaceCommand{
			Profile: profileFromLegacy(*legacy),
		})
		return replaceErr
	})
}

func mapSiteProfileError(err error) error {
	var conflict *siteprofile.RevisionConflictError
	var validation *siteprofile.ValidationError
	switch {
	case errors.As(err, &conflict):
		return naverr.RevisionConflict()
	case errors.As(err, &validation):
		return naverr.Validation("profile", "invalid", map[string]any{"message": validation.Error()})
	default:
		return err
	}
}

func (s *Service) catalog(ctx context.Context, filter dao.LinkFilter) (*Catalog, error) {
	categories, err := s.store.Categories(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := s.store.Groups(ctx)
	if err != nil {
		return nil, err
	}
	links, err := s.store.Links(ctx, filter)
	if err != nil {
		return nil, err
	}
	settings, settingsErr := s.PublicSite(ctx)
	if settingsErr != nil {
		return nil, settingsErr
	}
	return &Catalog{Site: settings, Categories: categories, Groups: groups, Links: links}, nil
}

func siteFromSnapshot(snapshot siteprofile.Snapshot, settings model.SiteSettings) Site {
	runtimeRevision := normalizedRuntimeRevision(settings.RuntimeRevision)
	return Site{
		Revision: uint64(snapshot.Revision), RuntimeRevision: runtimeRevision,
		ETag: consumerETag(snapshot.ETag, runtimeRevision, settings.SearchPlaceholder),
		Name: snapshot.Profile.Identity.Name, Title: snapshot.Profile.Identity.Tagline,
		Description: snapshot.Profile.Identity.Description, SearchPlaceholder: settings.SearchPlaceholder,
		FooterTagline: snapshot.Profile.Footer.Tagline,
	}
}

func normalizedRuntimeRevision(value uint64) uint64 {
	if value == 0 {
		return 1
	}
	return value
}

func consumerETag(profileETag string, runtimeRevision uint64, searchPlaceholder string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf(
		"%s:%d:%s", profileETag, runtimeRevision, strings.TrimSpace(searchPlaceholder),
	)))
	return fmt.Sprintf(`"nav-settings-r%d-%s"`, runtimeRevision, hex.EncodeToString(sum[:8]))
}

func profileFromLegacy(settings model.LegacySiteSettings) siteprofile.Profile {
	return siteprofile.Profile{
		Identity: siteprofile.Identity{
			Name: strings.TrimSpace(settings.Name), Tagline: strings.TrimSpace(settings.Title),
			Description: strings.TrimSpace(settings.Description),
		},
		Branding:     siteprofile.Branding{},
		Announcement: siteprofile.Announcement{},
		Support:      siteprofile.Support{Contacts: []siteprofile.Contact{}},
		Footer: siteprofile.Footer{
			Tagline:    strings.TrimSpace(settings.FooterTagline),
			LinkGroups: []siteprofile.LinkGroup{}, Social: []siteprofile.SocialLink{},
			Legal: []siteprofile.Link{}, Compliance: siteprofile.Compliance{Records: []siteprofile.ComplianceRecord{}},
		},
	}
}

func validateCategory(input model.Category) (*model.Category, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Icon = strings.TrimSpace(input.Icon)
	if input.Title == "" {
		return nil, naverr.Validation("title", "required", nil)
	}
	if input.Icon == "" {
		input.Icon = "i-tabler-folder"
	}
	if input.SortOrder < 0 {
		return nil, naverr.Validation("sortOrder", "minimum", map[string]any{"min": 0})
	}
	return &input, nil
}

func (s *Service) validateGroup(ctx context.Context, input model.Group) (*model.Group, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.CategoryID = strings.TrimSpace(input.CategoryID)
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	if input.CategoryID == "" {
		return nil, naverr.Validation("categoryId", "required", nil)
	}
	if input.Title == "" {
		return nil, naverr.Validation("title", "required", nil)
	}
	if input.SortOrder < 0 {
		return nil, naverr.Validation("sortOrder", "minimum", map[string]any{"min": 0})
	}
	categories, err := s.store.Categories(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(categories, func(category *model.Category) bool { return category.ID == input.CategoryID }) {
		return nil, naverr.Validation("categoryId", "not_found", nil)
	}
	return &input, nil
}

func (s *Service) validate(ctx context.Context, input LinkInput) (*model.Link, error) {
	input.ID = strings.TrimSpace(input.ID)
	input.CategoryID = strings.TrimSpace(input.CategoryID)
	input.GroupID = strings.TrimSpace(input.GroupID)
	input.Title = strings.TrimSpace(input.Title)
	input.URL = strings.TrimSpace(input.URL)
	input.Description = strings.TrimSpace(input.Description)
	input.Kind = strings.TrimSpace(input.Kind)
	input.Status = strings.TrimSpace(input.Status)
	input.SubmitterSub = strings.TrimSpace(input.SubmitterSub)
	if input.Title == "" {
		return nil, naverr.Validation("title", "required", nil)
	}
	if input.URL == "" {
		return nil, naverr.Validation("url", "required", nil)
	}
	if input.Description == "" {
		return nil, naverr.Validation("description", "required", nil)
	}
	parsed, err := url.ParseRequestURI(input.URL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, naverr.Validation("url", "absolute_http_url", nil)
	}
	if !slices.Contains(allowedKinds, input.Kind) {
		return nil, naverr.Validation("kind", "one_of", map[string]any{"allowed": allowedKinds})
	}
	if input.Status == "" {
		input.Status = "draft"
	}
	if !slices.Contains(allowedStatuses, input.Status) {
		return nil, naverr.Validation("status", "one_of", map[string]any{"allowed": allowedStatuses})
	}
	if input.SortOrder < 0 {
		return nil, naverr.Validation("sortOrder", "minimum", map[string]any{"min": 0})
	}
	validGroup, err := s.store.GroupBelongsToCategory(ctx, input.GroupID, input.CategoryID)
	if err != nil {
		return nil, err
	}
	if !validGroup {
		return nil, naverr.Validation("groupId", "category_mismatch", map[string]any{"categoryId": input.CategoryID})
	}
	return &model.Link{
		ID:           input.ID,
		CategoryID:   input.CategoryID,
		GroupID:      input.GroupID,
		Title:        input.Title,
		URL:          input.URL,
		Description:  input.Description,
		Tags:         normalize(input.Tags, 6),
		Keywords:     normalize(input.Keywords, 12),
		Kind:         input.Kind,
		Featured:     input.Featured,
		Status:       input.Status,
		SortOrder:    input.SortOrder,
		SubmitterSub: input.SubmitterSub,
	}, nil
}

func normalize(values []string, limit int) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, min(len(values), limit))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, value)
		if len(out) == limit {
			break
		}
	}
	return out
}
