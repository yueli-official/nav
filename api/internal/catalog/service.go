package catalog

import (
	"context"
	"net/url"
	"slices"
	"strings"

	"github.com/google/uuid"

	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/model"
	"platform/products/nav/api/internal/naverr"
)

const StatusPublished = "published"

var allowedKinds = []string{"official", "tool", "community", "learning", "resource", "reference", "research"}
var allowedStatuses = []string{StatusPublished, "draft", "archived"}

type Store interface {
	Categories(context.Context) ([]*model.Category, error)
	Groups(context.Context) ([]*model.Group, error)
	Links(context.Context, dao.LinkFilter) ([]*model.Link, error)
	GroupBelongsToCategory(context.Context, string, string) (bool, error)
	LinkExists(context.Context, string) (bool, error)
	InsertLink(context.Context, *model.Link) error
	UpdateLink(context.Context, *model.Link) (bool, error)
	DeleteLink(context.Context, string) (bool, error)
}

type Site struct {
	Name              string
	Title             string
	Description       string
	SearchPlaceholder string
}

type LinkInput struct {
	ID          string
	CategoryID  string
	GroupID     string
	Title       string
	URL         string
	Description string
	Tags        []string
	Keywords    []string
	Kind        string
	Featured    bool
	Status      string
	SortOrder   int
}

type Catalog struct {
	Site       Site
	Categories []*model.Category
	Groups     []*model.Group
	Links      []*model.Link
}

type Service struct {
	store Store
	site  Site
}

func New(store Store, site Site) *Service {
	return &Service{store: store, site: site}
}

func (s *Service) PublicCatalog(ctx context.Context) (*Catalog, error) {
	return s.catalog(ctx, dao.LinkFilter{Status: StatusPublished})
}

func (s *Service) AdminLinks(ctx context.Context, filter dao.LinkFilter) ([]*model.Link, error) {
	return s.store.Links(ctx, filter)
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
	if err := s.store.InsertLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *Service) UpdateLink(ctx context.Context, id string, input LinkInput) (*model.Link, error) {
	input.ID = strings.TrimSpace(id)
	link, err := s.validate(ctx, input)
	if err != nil {
		return nil, err
	}
	updated, err := s.store.UpdateLink(ctx, link)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, naverr.NotFound(id)
	}
	return link, nil
}

func (s *Service) DeleteLink(ctx context.Context, id string) error {
	deleted, err := s.store.DeleteLink(ctx, strings.TrimSpace(id))
	if err != nil {
		return err
	}
	if !deleted {
		return naverr.NotFound(id)
	}
	return nil
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
	return &Catalog{Site: s.site, Categories: categories, Groups: groups, Links: links}, nil
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
		ID:          input.ID,
		CategoryID:  input.CategoryID,
		GroupID:     input.GroupID,
		Title:       input.Title,
		URL:         input.URL,
		Description: input.Description,
		Tags:        normalize(input.Tags, 6),
		Keywords:    normalize(input.Keywords, 12),
		Kind:        input.Kind,
		Featured:    input.Featured,
		Status:      input.Status,
		SortOrder:   input.SortOrder,
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
