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
	Health     string
	Sort       string
	Direction  string
	Page       int
	Size       int
}

type linkHealthMutation struct {
	HealthStatus     string      `orm:"health_status"`
	HealthHTTPStatus any         `orm:"health_http_status"`
	HealthLatencyMS  any         `orm:"health_latency_ms"`
	HealthError      string      `orm:"health_error"`
	LastCheckedAt    *gtime.Time `orm:"last_checked_at"`
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
	SearchPlaceholder string      `orm:"search_placeholder"`
	RuntimeRevision   uint64      `orm:"runtime_revision"`
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
	query = query.Order(linkOrder(filter.Sort, filter.Direction))
	err := query.Scan(&links)
	return links, gerror.Wrap(err, "list navigation links")
}

func (p *PG) LinksByIDs(ctx context.Context, ids []string) ([]*model.Link, error) {
	if len(ids) == 0 {
		return []*model.Link{}, nil
	}
	var links []*model.Link
	err := p.db.Model(tLinks).Ctx(ctx).WhereIn("id", ids).OrderAsc("title").Scan(&links)
	return links, gerror.Wrap(err, "list navigation links by ids")
}

func (p *PG) LinkByID(ctx context.Context, id string) (*model.Link, error) {
	record, err := p.db.Model(tLinks).Ctx(ctx).Where("id", id).One()
	if err != nil {
		return nil, gerror.Wrap(err, "get navigation link")
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var link model.Link
	if err := record.Struct(&link); err != nil {
		return nil, gerror.Wrap(err, "scan navigation link")
	}
	return &link, nil
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
	switch filter.Health {
	case "unchecked":
		query = query.Where("health_status", "unchecked")
	case "healthy", "redirected", "broken", "timeout", "error":
		query = query.Where("health_status", filter.Health)
	case "issue":
		query = query.WhereIn("health_status", []string{"redirected", "broken", "timeout", "error"})
	}
	return query
}

func (p *PG) LinkHealthCounts(ctx context.Context) (map[string]int, error) {
	counts := map[string]int{"all": 0, "unchecked": 0, "healthy": 0, "redirected": 0, "broken": 0, "timeout": 0, "error": 0}
	rows, err := p.db.Model(tLinks).Ctx(ctx).Fields("health_status, COUNT(*) AS total").Group("health_status").All()
	if err != nil {
		return nil, gerror.Wrap(err, "count navigation link health")
	}
	for _, row := range rows {
		status, count := row["health_status"].String(), row["total"].Int()
		counts[status] = count
		counts["all"] += count
	}
	return counts, nil
}

func (p *PG) RecordClick(ctx context.Context, id string) (bool, error) {
	result, err := p.db.Model(tLinks).Ctx(ctx).Where("id", id).Where("status", "published").Data(struct {
		ClickCount    any         `orm:"click_count"`
		LastClickedAt *gtime.Time `orm:"last_clicked_at"`
	}{ClickCount: gdb.Raw("click_count + 1"), LastClickedAt: gtime.Now()}).Update()
	if err != nil {
		return false, gerror.Wrap(err, "record navigation link click")
	}
	changed, err := result.RowsAffected()
	return changed > 0, gerror.Wrap(err, "read navigation click result")
}

func (p *PG) UpdateLinkHealth(ctx context.Context, id string, health model.LinkHealth) (bool, error) {
	var httpStatus any
	if health.HTTPStatus > 0 {
		httpStatus = health.HTTPStatus
	} else {
		httpStatus = gdb.Raw("NULL")
	}
	var latency any
	if health.LatencyMS > 0 {
		latency = health.LatencyMS
	} else {
		latency = gdb.Raw("NULL")
	}
	result, err := p.db.Model(tLinks).Ctx(ctx).Where("id", id).Data(linkHealthMutation{
		HealthStatus: health.Status, HealthHTTPStatus: httpStatus, HealthLatencyMS: latency,
		HealthError: health.Error, LastCheckedAt: health.CheckedAt,
	}).Update()
	if err != nil {
		return false, gerror.Wrap(err, "update navigation link health")
	}
	changed, err := result.RowsAffected()
	return changed > 0, gerror.Wrap(err, "read navigation health update result")
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
	return p.RenameTagWithHook(ctx, source, target, nil)
}

func (p *PG) RenameTagWithHook(ctx context.Context, source, target string, hook TransactionHook) (int, error) {
	needle, _ := json.Marshal([]string{source})
	changed := 0
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Ctx(ctx).Exec(`
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
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		changed = int(affected)
		return runTransactionHook(ctx, tx, hook)
	})
	return changed, gerror.Wrap(err, "rename navigation tag")
}

func (p *PG) DeleteTag(ctx context.Context, name string) (int, error) {
	return p.DeleteTagWithHook(ctx, name, nil)
}

func (p *PG) DeleteTagWithHook(ctx context.Context, name string, hook TransactionHook) (int, error) {
	needle, _ := json.Marshal([]string{name})
	changed := 0
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Ctx(ctx).Exec(`
		UPDATE nav_links
		SET tags = (
			SELECT COALESCE(jsonb_agg(tag), '[]'::jsonb)
			FROM jsonb_array_elements_text(nav_links.tags) AS tag
			WHERE tag <> ?
		), updated_at = NOW()
		WHERE tags @> ?::jsonb`, name, string(needle))
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		changed = int(affected)
		return runTransactionHook(ctx, tx, hook)
	})
	return changed, gerror.Wrap(err, "delete navigation tag")
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
	return p.InsertLinkWithHook(ctx, link, nil)
}

func (p *PG) InsertLinkWithHook(ctx context.Context, link *model.Link, hook TransactionHook) error {
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(tLinks).Ctx(ctx).Data(insertMutation(link)).Insert(); err != nil {
			return err
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return gerror.Wrap(err, "insert navigation link")
}

func (p *PG) UpdateLink(ctx context.Context, link *model.Link) (bool, error) {
	return p.UpdateLinkWithHook(ctx, link, nil)
}

func (p *PG) UpdateLinkWithHook(ctx context.Context, link *model.Link, hook TransactionHook) (bool, error) {
	updated := false
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(tLinks).Ctx(ctx).Where("id", link.ID).Data(mutation(link)).Update()
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		updated = affected > 0
		if !updated {
			return nil
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return updated, gerror.Wrap(err, "update navigation link")
}

func (p *PG) DeleteLink(ctx context.Context, id string) (bool, error) {
	return p.DeleteLinkWithHook(ctx, id, nil)
}

func (p *PG) DeleteLinkWithHook(ctx context.Context, id string, hook TransactionHook) (bool, error) {
	deleted := false
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(tLinks).Ctx(ctx).Where("id", id).Delete()
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		deleted = affected > 0
		if !deleted {
			return nil
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return deleted, gerror.Wrap(err, "delete navigation link")
}

func (p *PG) BulkUpdateLinks(ctx context.Context, ids []string, status string) (int, error) {
	return p.BulkUpdateLinksWithHook(ctx, ids, status, nil)
}

func (p *PG) BulkUpdateLinksWithHook(ctx context.Context, ids []string, status string, hook TransactionHook) (int, error) {
	changed := 0
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(tLinks).Ctx(ctx).WhereIn("id", ids).Data(linkStatusMutation{Status: status, UpdatedAt: gtime.Now()}).Update()
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		changed = int(affected)
		return runTransactionHook(ctx, tx, hook)
	})
	return changed, gerror.Wrap(err, "bulk update navigation links")
}

func (p *PG) BulkDeleteLinks(ctx context.Context, ids []string) (int, error) {
	return p.BulkDeleteLinksWithHook(ctx, ids, nil)
}

func (p *PG) BulkDeleteLinksWithHook(ctx context.Context, ids []string, hook TransactionHook) (int, error) {
	changed := 0
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(tLinks).Ctx(ctx).WhereIn("id", ids).Delete()
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		changed = int(affected)
		return runTransactionHook(ctx, tx, hook)
	})
	return changed, gerror.Wrap(err, "bulk delete navigation links")
}

type linkStatusMutation struct {
	Status    string      `orm:"status"`
	UpdatedAt *gtime.Time `orm:"updated_at"`
}

func (p *PG) InsertCategory(ctx context.Context, category *model.Category) error {
	return p.InsertCategoryWithHook(ctx, category, nil)
}

func (p *PG) InsertCategoryWithHook(ctx context.Context, category *model.Category, hook TransactionHook) error {
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(tCategories).Ctx(ctx).Data(categoryInsertMutation{ID: category.ID, Title: category.Title, Description: category.Description, Icon: category.Icon, SortOrder: category.SortOrder}).Insert(); err != nil {
			return err
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return gerror.Wrap(err, "insert navigation category")
}

func (p *PG) UpdateCategory(ctx context.Context, category *model.Category) (bool, error) {
	return p.UpdateCategoryWithHook(ctx, category, nil)
}

func (p *PG) UpdateCategoryWithHook(ctx context.Context, category *model.Category, hook TransactionHook) (bool, error) {
	updated := false
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, err := tx.Model(tCategories).Ctx(ctx).Where("id", category.ID).Data(categoryData(category)).Update()
		if err != nil {
			return err
		}
		changed, err := result.RowsAffected()
		if err != nil {
			return err
		}
		updated = changed > 0
		if !updated {
			return nil
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return updated, gerror.Wrap(err, "update navigation category")
}

func (p *PG) DeleteCategory(ctx context.Context, id string) (bool, error) {
	return p.DeleteCategoryWithHook(ctx, id, nil)
}

func (p *PG) DeleteCategoryWithHook(ctx context.Context, id string, hook TransactionHook) (bool, error) {
	deleted := false
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		count, err := tx.Model(tGroups).Ctx(ctx).Where("category_id", id).Count()
		if err != nil || count > 0 {
			return err
		}
		result, err := tx.Model(tCategories).Ctx(ctx).Where("id", id).Delete()
		if err != nil {
			return err
		}
		changed, err := result.RowsAffected()
		if err != nil {
			return err
		}
		deleted = changed > 0
		if !deleted {
			return nil
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return deleted, gerror.Wrap(err, "delete navigation category")
}

func (p *PG) InsertGroup(ctx context.Context, group *model.Group) error {
	return p.InsertGroupWithHook(ctx, group, nil)
}

func (p *PG) InsertGroupWithHook(ctx context.Context, group *model.Group, hook TransactionHook) error {
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(tGroups).Ctx(ctx).Data(groupInsertMutation{ID: group.ID, CategoryID: group.CategoryID, Title: group.Title, Description: group.Description, SortOrder: group.SortOrder}).Insert(); err != nil {
			return err
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return gerror.Wrap(err, "insert navigation group")
}

func (p *PG) UpdateGroup(ctx context.Context, group *model.Group) (bool, error) {
	return p.UpdateGroupWithHook(ctx, group, nil)
}

func (p *PG) UpdateGroupWithHook(ctx context.Context, group *model.Group, hook TransactionHook) (bool, error) {
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
		if updateErr != nil {
			return updateErr
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return err == nil, gerror.Wrap(err, "update navigation group")
}

func (p *PG) DeleteGroup(ctx context.Context, id string) (bool, error) {
	return p.DeleteGroupWithHook(ctx, id, nil)
}

func (p *PG) DeleteGroupWithHook(ctx context.Context, id string, hook TransactionHook) (bool, error) {
	deleted := false
	err := p.db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		count, err := tx.Model(tLinks).Ctx(ctx).Where("group_id", id).Count()
		if err != nil || count > 0 {
			return err
		}
		result, err := tx.Model(tGroups).Ctx(ctx).Where("id", id).Delete()
		if err != nil {
			return err
		}
		changed, err := result.RowsAffected()
		if err != nil {
			return err
		}
		deleted = changed > 0
		if !deleted {
			return nil
		}
		return runTransactionHook(ctx, tx, hook)
	})
	return deleted, gerror.Wrap(err, "delete navigation group")
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

func (p *PG) LegacySiteSettings(ctx context.Context) (*model.LegacySiteSettings, error) {
	record, err := p.db.Model(tSettings).Ctx(ctx).Where("id", 1).One()
	if err != nil {
		return nil, gerror.Wrap(err, "get legacy navigation site settings")
	}
	if record.IsEmpty() {
		return nil, nil
	}
	var settings model.LegacySiteSettings
	if err := record.Struct(&settings); err != nil {
		return nil, gerror.Wrap(err, "scan legacy navigation site settings")
	}
	return &settings, nil
}

func (p *PG) UpsertSiteSettings(ctx context.Context, settings *model.SiteSettings) error {
	data := siteSettingsMutation{
		ID: 1, SearchPlaceholder: settings.SearchPlaceholder,
		RuntimeRevision: settings.RuntimeRevision, UpdatedAt: gtime.Now(),
	}
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

func linkOrder(sort, direction string) string {
	orderDirection := "ASC"
	if direction == "desc" {
		orderDirection = "DESC"
	}
	switch sort {
	case "popular":
		return "click_count DESC, featured DESC, title ASC, id ASC"
	case "health":
		return "last_checked_at ASC, title ASC, id ASC"
	case "updated":
		return "updated_at " + orderDirection + ", title ASC, id ASC"
	case "title":
		return "title " + orderDirection + ", id ASC"
	case "published":
		return "published_at " + orderDirection + " NULLS LAST, title ASC, id ASC"
	default:
		return "featured DESC, sort_order " + orderDirection + ", title " + orderDirection + ", id ASC"
	}
}
