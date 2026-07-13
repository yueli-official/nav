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
	tSettings   = "nav_site_settings"
)

type LinkFilter struct {
	Query      string
	CategoryID string
	GroupID    string
	Status     string
	Tag        string
	Page       int
	Size       int
}

type categoryMutation struct {
	Title       string `orm:"title"`
	Description string `orm:"description"`
	Icon        string `orm:"icon"`
	SortOrder   int    `orm:"sort_order"`
}

type categoryInsertMutation struct {
	ID          string `orm:"id"`
	Title       string `orm:"title"`
	Description string `orm:"description"`
	Icon        string `orm:"icon"`
	SortOrder   int    `orm:"sort_order"`
}

type groupMutation struct {
	CategoryID  string `orm:"category_id"`
	Title       string `orm:"title"`
	Description string `orm:"description"`
	SortOrder   int    `orm:"sort_order"`
}

type groupInsertMutation struct {
	ID          string `orm:"id"`
	CategoryID  string `orm:"category_id"`
	Title       string `orm:"title"`
	Description string `orm:"description"`
	SortOrder   int    `orm:"sort_order"`
}

type siteSettingsMutation struct {
	ID                int         `orm:"id"`
	Name              string      `orm:"name"`
	Title             string      `orm:"title"`
	Description       string      `orm:"description"`
	SearchPlaceholder string      `orm:"search_placeholder"`
	FooterTagline     string      `orm:"footer_tagline"`
	UpdatedAt         *gtime.Time `orm:"updated_at"`
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
	query := p.linkQuery(ctx, filter)
	if filter.Size > 0 {
		page := max(filter.Page, 1)
		query = query.Limit((page-1)*filter.Size, filter.Size)
	}
	var links []*model.Link
	err := query.OrderDesc("featured").OrderAsc("sort_order").OrderAsc("title").Scan(&links)
	return links, gerror.Wrap(err, "list navigation links")
}

func (p *PG) CountLinks(ctx context.Context, filter LinkFilter) (int, error) {
	filter.Page = 0
	filter.Size = 0
	count, err := p.linkQuery(ctx, filter).Count()
	return count, gerror.Wrap(err, "count navigation links")
}

func (p *PG) LinkStatusCounts(ctx context.Context) (map[string]int, error) {
	counts := map[string]int{}
	all, err := p.db.Model(tLinks).Ctx(ctx).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "count all navigation links")
	}
	counts["all"] = all
	for _, status := range []string{"published", "draft", "archived"} {
		count, countErr := p.db.Model(tLinks).Ctx(ctx).Where("status", status).Count()
		if countErr != nil {
			return nil, gerror.Wrap(countErr, "count navigation links by status")
		}
		counts[status] = count
	}
	return counts, nil
}

func (p *PG) linkQuery(ctx context.Context, filter LinkFilter) *gdb.Model {
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
	if strings.TrimSpace(filter.Tag) != "" {
		tag, _ := json.Marshal([]string{strings.TrimSpace(filter.Tag)})
		query = query.Where("tags @> ?::jsonb", string(tag))
	}
	if strings.TrimSpace(filter.Query) != "" {
		// TODO(nav-search): the local/prod PG images must expose zhparser before
		// this small admin catalog can move to the platform tsvector + GIN path.
		// Until then, keep the documented development fallback explicit.
		like := "%" + strings.TrimSpace(filter.Query) + "%"
		query = query.Where("(title ILIKE ? OR url ILIKE ? OR description ILIKE ?)", like, like, like)
	}
	return query
}

func (p *PG) Tags(ctx context.Context, query string) ([]*model.Tag, error) {
	args := []any{}
	where := ""
	if strings.TrimSpace(query) != "" {
		where = "WHERE tag ILIKE ?"
		args = append(args, "%"+strings.TrimSpace(query)+"%")
	}
	rows, err := p.db.Ctx(ctx).Raw(`
		SELECT tag AS name, COUNT(*)::int AS link_count
		FROM nav_links CROSS JOIN LATERAL jsonb_array_elements_text(tags) AS tag
		`+where+`
		GROUP BY tag
		ORDER BY COUNT(*) DESC, tag ASC`, args...).All()
	if err != nil {
		return nil, gerror.Wrap(err, "list navigation tags")
	}
	tags := make([]*model.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, &model.Tag{Name: row["name"].String(), LinkCount: row["link_count"].Int()})
	}
	return tags, nil
}

func (p *PG) RenameTag(ctx context.Context, source, target string) (int, error) {
	needle, _ := json.Marshal([]string{source})
	result, err := p.db.Exec(ctx, `
		UPDATE nav_links
		SET tags = (
			SELECT COALESCE(jsonb_agg(value), '[]'::jsonb)
			FROM (
				SELECT DISTINCT CASE WHEN tag = ? THEN ? ELSE tag END AS value
				FROM jsonb_array_elements_text(nav_links.tags) AS tag
			) normalized
		), updated_at = NOW()
		WHERE tags @> ?::jsonb`, source, target, string(needle))
	if err != nil {
		return 0, gerror.Wrap(err, "rename navigation tag")
	}
	changed, err := result.RowsAffected()
	return int(changed), gerror.Wrap(err, "read navigation tag rename result")
}

func (p *PG) DeleteTag(ctx context.Context, name string) (int, error) {
	needle, _ := json.Marshal([]string{name})
	result, err := p.db.Exec(ctx, `
		UPDATE nav_links
		SET tags = (
			SELECT COALESCE(jsonb_agg(tag), '[]'::jsonb)
			FROM jsonb_array_elements_text(nav_links.tags) AS tag
			WHERE tag <> ?
		), updated_at = NOW()
		WHERE tags @> ?::jsonb`, name, string(needle))
	if err != nil {
		return 0, gerror.Wrap(err, "delete navigation tag")
	}
	changed, err := result.RowsAffected()
	return int(changed), gerror.Wrap(err, "read navigation tag delete result")
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

func (p *PG) BulkUpdateLinks(ctx context.Context, ids []string, status string) (int, error) {
	result, err := p.db.Model(tLinks).Ctx(ctx).WhereIn("id", ids).Data(linkStatusMutation{Status: status, UpdatedAt: gtime.Now()}).Update()
	if err != nil {
		return 0, gerror.Wrap(err, "bulk update navigation links")
	}
	changed, err := result.RowsAffected()
	return int(changed), gerror.Wrap(err, "read navigation bulk update result")
}

func (p *PG) BulkDeleteLinks(ctx context.Context, ids []string) (int, error) {
	result, err := p.db.Model(tLinks).Ctx(ctx).WhereIn("id", ids).Delete()
	if err != nil {
		return 0, gerror.Wrap(err, "bulk delete navigation links")
	}
	changed, err := result.RowsAffected()
	return int(changed), gerror.Wrap(err, "read navigation bulk delete result")
}

type linkStatusMutation struct {
	Status    string      `orm:"status"`
	UpdatedAt *gtime.Time `orm:"updated_at"`
}

func (p *PG) InsertCategory(ctx context.Context, category *model.Category) error {
	_, err := p.db.Model(tCategories).Ctx(ctx).Data(categoryInsertMutation{ID: category.ID, Title: category.Title, Description: category.Description, Icon: category.Icon, SortOrder: category.SortOrder}).Insert()
	return gerror.Wrap(err, "insert navigation category")
}

func (p *PG) UpdateCategory(ctx context.Context, category *model.Category) (bool, error) {
	result, err := p.db.Model(tCategories).Ctx(ctx).Where("id", category.ID).Data(categoryData(category)).Update()
	if err != nil {
		return false, gerror.Wrap(err, "update navigation category")
	}
	changed, err := result.RowsAffected()
	return changed > 0, gerror.Wrap(err, "read navigation category update result")
}

func (p *PG) DeleteCategory(ctx context.Context, id string) (bool, error) {
	count, err := p.db.Model(tGroups).Ctx(ctx).Where("category_id", id).Count()
	if err != nil {
		return false, gerror.Wrap(err, "count navigation category groups")
	}
	if count > 0 {
		return false, nil
	}
	result, err := p.db.Model(tCategories).Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return false, gerror.Wrap(err, "delete navigation category")
	}
	changed, err := result.RowsAffected()
	return changed > 0, gerror.Wrap(err, "read navigation category delete result")
}

func (p *PG) InsertGroup(ctx context.Context, group *model.Group) error {
	_, err := p.db.Model(tGroups).Ctx(ctx).Data(groupInsertMutation{ID: group.ID, CategoryID: group.CategoryID, Title: group.Title, Description: group.Description, SortOrder: group.SortOrder}).Insert()
	return gerror.Wrap(err, "insert navigation group")
}

func (p *PG) UpdateGroup(ctx context.Context, group *model.Group) (bool, error) {
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, updateErr := tx.Model(tGroups).Ctx(ctx).Where("id", group.ID).Data(groupData(group)).Update()
		if updateErr != nil {
			return updateErr
		}
		changed, affectedErr := result.RowsAffected()
		if affectedErr != nil {
			return affectedErr
		}
		if changed == 0 {
			return gerror.New("navigation group not found")
		}
		_, updateErr = tx.Model(tLinks).Ctx(ctx).Where("group_id", group.ID).Data(struct {
			CategoryID string `orm:"category_id"`
		}{CategoryID: group.CategoryID}).Update()
		return updateErr
	})
	return err == nil, gerror.Wrap(err, "update navigation group")
}

func (p *PG) DeleteGroup(ctx context.Context, id string) (bool, error) {
	count, err := p.db.Model(tLinks).Ctx(ctx).Where("group_id", id).Count()
	if err != nil {
		return false, gerror.Wrap(err, "count navigation group links")
	}
	if count > 0 {
		return false, nil
	}
	result, err := p.db.Model(tGroups).Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return false, gerror.Wrap(err, "delete navigation group")
	}
	changed, err := result.RowsAffected()
	return changed > 0, gerror.Wrap(err, "read navigation group delete result")
}

func (p *PG) SiteSettings(ctx context.Context) (*model.SiteSettings, error) {
	record, err := p.db.Model(tSettings).Ctx(ctx).Where("id", 1).One()
	if err != nil {
		return nil, gerror.Wrap(err, "get navigation site settings")
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var settings model.SiteSettings
	if err := record.Struct(&settings); err != nil {
		return nil, gerror.Wrap(err, "scan navigation site settings")
	}
	return &settings, nil
}

func (p *PG) UpsertSiteSettings(ctx context.Context, settings *model.SiteSettings) error {
	data := siteSettingsMutation{ID: 1, Name: settings.Name, Title: settings.Title, Description: settings.Description, SearchPlaceholder: settings.SearchPlaceholder, FooterTagline: settings.FooterTagline, UpdatedAt: gtime.Now()}
	_, err := p.db.Model(tSettings).Ctx(ctx).Data(data).OnConflict("id").Save()
	return gerror.Wrap(err, "save navigation site settings")
}

func categoryData(category *model.Category) categoryMutation {
	return categoryMutation{Title: category.Title, Description: category.Description, Icon: category.Icon, SortOrder: category.SortOrder}
}

func groupData(group *model.Group) groupMutation {
	return groupMutation{CategoryID: group.CategoryID, Title: group.Title, Description: group.Description, SortOrder: group.SortOrder}
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
