package dao

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"platform/products/nav/api/internal/model"
)

const (
	tCategories = "nav_categories"
	tGroups     = "nav_groups"
	tLinks      = "nav_links"
)

type LinkFilter struct {
	Query      string
	CategoryID string
	GroupID    string
	Status     string
}

type linkMutation struct {
	CategoryID  string      `orm:"category_id"`
	GroupID     string      `orm:"group_id"`
	Title       string      `orm:"title"`
	URL         string      `orm:"url"`
	Description string      `orm:"description"`
	Tags        string      `orm:"tags"`
	Keywords    string      `orm:"keywords"`
	Kind        string      `orm:"kind"`
	Featured    bool        `orm:"featured"`
	Status      string      `orm:"status"`
	SortOrder   int         `orm:"sort_order"`
	UpdatedAt   *gtime.Time `orm:"updated_at,omitempty"`
}

type linkInsertMutation struct {
	ID          string      `orm:"id"`
	CategoryID  string      `orm:"category_id"`
	GroupID     string      `orm:"group_id"`
	Title       string      `orm:"title"`
	URL         string      `orm:"url"`
	Description string      `orm:"description"`
	Tags        string      `orm:"tags"`
	Keywords    string      `orm:"keywords"`
	Kind        string      `orm:"kind"`
	Featured    bool        `orm:"featured"`
	Status      string      `orm:"status"`
	SortOrder   int         `orm:"sort_order"`
	UpdatedAt   *gtime.Time `orm:"updated_at,omitempty"`
}

type PG struct {
	db gdb.DB
}

func NewPG(db gdb.DB) *PG {
	return &PG{db: db}
}

func (p *PG) Categories(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	err := p.db.Model(tCategories).Ctx(ctx).OrderAsc("sort_order").OrderAsc("title").Scan(&categories)
	return categories, gerror.Wrap(err, "list navigation categories")
}

func (p *PG) Groups(ctx context.Context) ([]*model.Group, error) {
	var groups []*model.Group
	err := p.db.Model(tGroups).Ctx(ctx).OrderAsc("sort_order").OrderAsc("title").Scan(&groups)
	return groups, gerror.Wrap(err, "list navigation groups")
}

func (p *PG) Links(ctx context.Context, filter LinkFilter) ([]*model.Link, error) {
	query := p.db.Model(tLinks).Ctx(ctx)
	if filter.Status != "" {
		query = query.Where("status", filter.Status)
	}
	if filter.CategoryID != "" {
		query = query.Where("category_id", filter.CategoryID)
	}
	if filter.GroupID != "" {
		query = query.Where("group_id", filter.GroupID)
	}
	if strings.TrimSpace(filter.Query) != "" {
		// TODO(nav-search): the local/prod PG images must expose zhparser before
		// this small admin catalog can move to the platform tsvector + GIN path.
		// Until then, keep the documented development fallback explicit.
		like := "%" + strings.TrimSpace(filter.Query) + "%"
		query = query.Where("(title ILIKE ? OR url ILIKE ? OR description ILIKE ?)", like, like, like)
	}
	var links []*model.Link
	err := query.OrderDesc("featured").OrderAsc("sort_order").OrderAsc("title").Scan(&links)
	return links, gerror.Wrap(err, "list navigation links")
}

func (p *PG) GroupBelongsToCategory(ctx context.Context, groupID, categoryID string) (bool, error) {
	count, err := p.db.Model(tGroups).Ctx(ctx).Where("id", groupID).Where("category_id", categoryID).Count()
	return count > 0, gerror.Wrap(err, "validate navigation group")
}

func (p *PG) LinkExists(ctx context.Context, id string) (bool, error) {
	count, err := p.db.Model(tLinks).Ctx(ctx).Where("id", id).Count()
	return count > 0, gerror.Wrap(err, "check navigation link")
}

func (p *PG) InsertLink(ctx context.Context, link *model.Link) error {
	_, err := p.db.Model(tLinks).Ctx(ctx).Data(insertMutation(link)).Insert()
	return gerror.Wrap(err, "insert navigation link")
}

func (p *PG) UpdateLink(ctx context.Context, link *model.Link) (bool, error) {
	result, err := p.db.Model(tLinks).Ctx(ctx).Where("id", link.ID).Data(mutation(link)).Update()
	if err != nil {
		return false, gerror.Wrap(err, "update navigation link")
	}
	affected, err := result.RowsAffected()
	return affected > 0, gerror.Wrap(err, "read navigation update result")
}

func (p *PG) DeleteLink(ctx context.Context, id string) (bool, error) {
	result, err := p.db.Model(tLinks).Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return false, gerror.Wrap(err, "delete navigation link")
	}
	affected, err := result.RowsAffected()
	return affected > 0, gerror.Wrap(err, "read navigation delete result")
}

func mutation(link *model.Link) linkMutation {
	tags, _ := json.Marshal(link.Tags)
	keywords, _ := json.Marshal(link.Keywords)
	return linkMutation{
		CategoryID:  link.CategoryID,
		GroupID:     link.GroupID,
		Title:       link.Title,
		URL:         link.URL,
		Description: link.Description,
		Tags:        string(tags),
		Keywords:    string(keywords),
		Kind:        link.Kind,
		Featured:    link.Featured,
		Status:      link.Status,
		SortOrder:   link.SortOrder,
		UpdatedAt:   gtime.Now(),
	}
}

func insertMutation(link *model.Link) linkInsertMutation {
	data := mutation(link)
	return linkInsertMutation{
		ID: link.ID, CategoryID: data.CategoryID, GroupID: data.GroupID,
		Title: data.Title, URL: data.URL, Description: data.Description,
		Tags: data.Tags, Keywords: data.Keywords, Kind: data.Kind,
		Featured: data.Featured, Status: data.Status, SortOrder: data.SortOrder,
		UpdatedAt: data.UpdatedAt,
	}
}
