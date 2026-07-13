package controller

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	v1 "platform/products/nav/api/api/v1"
	"platform/products/nav/api/internal/catalog"
	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/model"
)

type Public struct {
	service *catalog.Service
}

func NewPublic(service *catalog.Service) *Public {
	return &Public{service: service}
}

func (c *Public) GetCatalog(ctx context.Context, _ *v1.GetCatalogReq) (*v1.GetCatalogRes, error) {
	value, err := c.service.PublicCatalog(ctx)
	if err != nil {
		return nil, err
	}
	return catalogResponse(value), nil
}

func (c *Public) GetGroup(ctx context.Context, req *v1.GetGroupReq) (*v1.GetGroupRes, error) {
	page, err := c.service.PublicGroup(ctx, req.GroupID, req.Page, req.Size, req.Sort)
	if err != nil {
		return nil, err
	}
	items := make([]v1.LinkView, 0, len(page.Links))
	for _, link := range page.Links {
		items = append(items, linkView(link, false))
	}
	category := categoryView(page.Category)
	group := groupView(page.Group)
	group.LinkCount = page.Total
	return &v1.GetGroupRes{
		Site: settingsView(&model.SiteSettings{
			Name: page.Site.Name, Title: page.Site.Title, Description: page.Site.Description,
			SearchPlaceholder: page.Site.SearchPlaceholder, FooterTagline: page.Site.FooterTagline,
		}),
		Category: category, Group: group, Items: items,
		Total: page.Total, Page: page.Page, Size: page.Size,
	}, nil
}

func (c *Public) RecordClick(ctx context.Context, req *v1.RecordClickReq) (*v1.RecordClickRes, error) {
	recorded, err := c.service.RecordClick(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return &v1.RecordClickRes{Recorded: recorded}, nil
}

func (c *Public) GetFavicon(ctx context.Context, req *v1.GetFaviconReq) (*v1.GetFaviconRes, error) {
	data, contentType, err := c.service.Favicon(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	request := ghttp.RequestFromCtx(ctx)
	request.Response.Header().Set("Content-Type", contentType)
	request.Response.Header().Set("Cache-Control", "public, max-age=86400, stale-while-revalidate=604800")
	request.Response.Header().Set("X-Content-Type-Options", "nosniff")
	request.Response.Write(data)
	return nil, nil
}

type Admin struct {
	service *catalog.Service
}

func NewAdmin(service *catalog.Service) *Admin {
	return &Admin{service: service}
}

func (c *Admin) AdminListLinks(ctx context.Context, req *v1.AdminListLinksReq) (*v1.AdminListLinksRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	page, err := c.service.AdminLinks(ctx, dao.LinkFilter{
		Query: req.Q, CategoryID: req.CategoryID, GroupID: req.GroupID, Status: req.Status,
		Tag: req.Tag, Sort: req.Sort, Direction: req.Direction, Page: req.Page, Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	structure, err := c.service.AdminStructure(ctx)
	if err != nil {
		return nil, err
	}
	views := make([]v1.LinkView, 0, len(page.Links))
	for _, link := range page.Links {
		views = append(views, linkView(link, true))
	}
	tags := make([]v1.TagView, 0, len(page.Tags))
	for _, tag := range page.Tags {
		tags = append(tags, v1.TagView{Name: tag.Name, LinkCount: tag.LinkCount})
	}
	return &v1.AdminListLinksRes{
		Links: views, Categories: categoryViews(structure, false), Tags: tags,
		Counts: v1.LifecycleCountsView{All: page.Counts["all"], Published: page.Counts["published"], Draft: page.Counts["draft"], Archived: page.Counts["archived"]},
		Total:  page.Total, Page: req.Page, Size: req.Size,
	}, nil
}

func (c *Admin) AdminListChecks(ctx context.Context, req *v1.AdminListChecksReq) (*v1.AdminListChecksRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	page, err := c.service.AdminChecks(ctx, dao.LinkFilter{Query: req.Q, Health: req.Health, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	links := make([]v1.LinkView, 0, len(page.Links))
	for _, link := range page.Links {
		links = append(links, linkView(link, true))
	}
	return &v1.AdminListChecksRes{
		Links: links,
		Counts: v1.HealthCountsView{
			All: page.Counts["all"], Unchecked: page.Counts["unchecked"], Healthy: page.Counts["healthy"],
			Redirected: page.Counts["redirected"], Broken: page.Counts["broken"],
			Timeout: page.Counts["timeout"], Error: page.Counts["error"],
		},
		Total: page.Total, Page: req.Page, Size: req.Size,
	}, nil
}

func (c *Admin) AdminRunChecks(ctx context.Context, req *v1.AdminRunChecksReq) (*v1.AdminRunChecksRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	results, err := c.service.RunChecks(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	views := make([]v1.LinkView, 0, len(results))
	for _, link := range results {
		views = append(views, linkView(link, true))
	}
	return &v1.AdminRunChecksRes{Checked: len(views), Results: views}, nil
}

func (c *Admin) AdminBulkLinks(ctx context.Context, req *v1.AdminBulkLinksReq) (*v1.AdminBulkLinksRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	result, err := c.service.BulkLinks(ctx, req.IDs, req.Action)
	if err != nil {
		return nil, err
	}
	return &v1.AdminBulkLinksRes{Changed: result.Changed, FailedIDs: result.FailedIDs}, nil
}

func (c *Admin) AdminListStructure(ctx context.Context, _ *v1.AdminListStructureReq) (*v1.AdminListStructureRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	structure, err := c.service.AdminStructure(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.AdminListStructureRes{Categories: categoryViews(structure, false)}, nil
}

func (c *Admin) AdminCreateCategory(ctx context.Context, req *v1.AdminCreateCategoryReq) (*v1.AdminCreateCategoryRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	category, err := c.service.CreateCategory(ctx, categoryInput(req.CategoryInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminCreateCategoryRes{Category: categoryView(category)}, nil
}

func (c *Admin) AdminUpdateCategory(ctx context.Context, req *v1.AdminUpdateCategoryReq) (*v1.AdminUpdateCategoryRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	category, err := c.service.UpdateCategory(ctx, req.ID, categoryInput(req.CategoryInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminUpdateCategoryRes{Category: categoryView(category)}, nil
}

func (c *Admin) AdminDeleteCategory(ctx context.Context, req *v1.AdminDeleteCategoryReq) (*v1.AdminDeleteCategoryRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := c.service.DeleteCategory(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.AdminDeleteCategoryRes{Deleted: true}, nil
}

func (c *Admin) AdminCreateGroup(ctx context.Context, req *v1.AdminCreateGroupReq) (*v1.AdminCreateGroupRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	group, err := c.service.CreateGroup(ctx, groupInput(req.GroupInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminCreateGroupRes{Group: groupView(group)}, nil
}

func (c *Admin) AdminUpdateGroup(ctx context.Context, req *v1.AdminUpdateGroupReq) (*v1.AdminUpdateGroupRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	group, err := c.service.UpdateGroup(ctx, req.ID, groupInput(req.GroupInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminUpdateGroupRes{Group: groupView(group)}, nil
}

func (c *Admin) AdminDeleteGroup(ctx context.Context, req *v1.AdminDeleteGroupReq) (*v1.AdminDeleteGroupRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := c.service.DeleteGroup(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.AdminDeleteGroupRes{Deleted: true}, nil
}

func (c *Admin) AdminListTags(ctx context.Context, req *v1.AdminListTagsReq) (*v1.AdminListTagsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	tags, err := c.service.Tags(ctx, req.Q)
	if err != nil {
		return nil, err
	}
	views := make([]v1.TagView, 0, len(tags))
	for _, tag := range tags {
		views = append(views, v1.TagView{Name: tag.Name, LinkCount: tag.LinkCount})
	}
	return &v1.AdminListTagsRes{Tags: views}, nil
}

func (c *Admin) AdminRenameTag(ctx context.Context, req *v1.AdminRenameTagReq) (*v1.AdminRenameTagRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	changed, err := c.service.RenameTag(ctx, req.Source, req.Target)
	if err != nil {
		return nil, err
	}
	return &v1.AdminRenameTagRes{Changed: changed}, nil
}

func (c *Admin) AdminDeleteTag(ctx context.Context, req *v1.AdminDeleteTagReq) (*v1.AdminDeleteTagRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	changed, err := c.service.DeleteTag(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &v1.AdminDeleteTagRes{Changed: changed}, nil
}

func (c *Admin) AdminGetSettings(ctx context.Context, _ *v1.AdminGetSettingsReq) (*v1.AdminGetSettingsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	settings, err := c.service.Settings(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.AdminGetSettingsRes{Settings: settingsView(settings)}, nil
}

func (c *Admin) AdminUpdateSettings(ctx context.Context, req *v1.AdminUpdateSettingsReq) (*v1.AdminUpdateSettingsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	settings, err := c.service.UpdateSettings(ctx, model.SiteSettings{Name: req.Name, Title: req.Title, Description: req.Description, SearchPlaceholder: req.SearchPlaceholder, FooterTagline: req.FooterTagline})
	if err != nil {
		return nil, err
	}
	return &v1.AdminUpdateSettingsRes{Settings: settingsView(settings)}, nil
}

func (c *Admin) AdminCreateLink(ctx context.Context, req *v1.AdminCreateLinkReq) (*v1.AdminCreateLinkRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	link, err := c.service.CreateLink(ctx, input(req.LinkInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminCreateLinkRes{Link: linkView(link, true)}, nil
}

func (c *Admin) AdminUpdateLink(ctx context.Context, req *v1.AdminUpdateLinkReq) (*v1.AdminUpdateLinkRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	link, err := c.service.UpdateLink(ctx, req.ID, input(req.LinkInput))
	if err != nil {
		return nil, err
	}
	return &v1.AdminUpdateLinkRes{Link: linkView(link, true)}, nil
}

func (c *Admin) AdminDeleteLink(ctx context.Context, req *v1.AdminDeleteLinkReq) (*v1.AdminDeleteLinkRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if err := c.service.DeleteLink(ctx, req.ID); err != nil {
		return nil, err
	}
	return &v1.AdminDeleteLinkRes{Deleted: true}, nil
}

func input(value v1.LinkInput) catalog.LinkInput {
	return catalog.LinkInput{
		CategoryID:  value.CategoryID,
		GroupID:     value.GroupID,
		Title:       value.Title,
		URL:         value.URL,
		Description: value.Description,
		Tags:        value.Tags,
		Keywords:    value.Keywords,
		Kind:        value.Kind,
		Featured:    value.Featured,
		Status:      value.Status,
		SortOrder:   value.SortOrder,
	}
}

func catalogResponse(value *catalog.Catalog) *v1.GetCatalogRes {
	return &v1.GetCatalogRes{
		Version: 1,
		Site: v1.SiteView{
			Name:              value.Site.Name,
			Title:             value.Site.Title,
			Description:       value.Site.Description,
			SearchPlaceholder: value.Site.SearchPlaceholder,
			FooterTagline:     value.Site.FooterTagline,
		},
		Categories: categoryViews(value, true),
		Stats: v1.StatsView{
			CategoryCount: len(value.Categories),
			GroupCount:    len(value.Groups),
			LinkCount:     len(value.Links),
		},
	}
}

func categoryInput(value v1.CategoryInput) model.Category {
	return model.Category{Title: value.Title, Description: value.Description, Icon: value.Icon, SortOrder: value.SortOrder}
}

func groupInput(value v1.GroupInput) model.Group {
	return model.Group{CategoryID: value.CategoryID, Title: value.Title, Description: value.Description, SortOrder: value.SortOrder}
}

func categoryView(category *model.Category) v1.CategoryView {
	return v1.CategoryView{ID: category.ID, Title: category.Title, Description: category.Description, Icon: category.Icon, SortOrder: category.SortOrder, Groups: []v1.GroupView{}}
}

func groupView(group *model.Group) v1.GroupView {
	return v1.GroupView{ID: group.ID, CategoryID: group.CategoryID, Title: group.Title, Description: group.Description, SortOrder: group.SortOrder, Items: []v1.LinkView{}}
}

func settingsView(settings *model.SiteSettings) v1.SiteView {
	return v1.SiteView{Name: settings.Name, Title: settings.Title, Description: settings.Description, SearchPlaceholder: settings.SearchPlaceholder, FooterTagline: settings.FooterTagline}
}

func categoryViews(value *catalog.Catalog, includeLinks bool) []v1.CategoryView {
	linksByGroup := make(map[string][]v1.LinkView, len(value.Groups))
	linkCountsByGroup := make(map[string]int, len(value.Groups))
	for _, link := range value.Links {
		linkCountsByGroup[link.GroupID]++
		if includeLinks {
			linksByGroup[link.GroupID] = append(linksByGroup[link.GroupID], linkView(link, false))
		}
	}
	groupsByCategory := make(map[string][]v1.GroupView, len(value.Categories))
	for _, group := range value.Groups {
		groupsByCategory[group.CategoryID] = append(groupsByCategory[group.CategoryID], v1.GroupView{
			ID:          group.ID,
			CategoryID:  group.CategoryID,
			Title:       group.Title,
			Description: group.Description,
			SortOrder:   group.SortOrder,
			LinkCount:   linkCountsByGroup[group.ID],
			Items:       linksByGroup[group.ID],
		})
	}
	views := make([]v1.CategoryView, 0, len(value.Categories))
	for _, category := range value.Categories {
		views = append(views, v1.CategoryView{
			ID:          category.ID,
			Title:       category.Title,
			Description: category.Description,
			Icon:        category.Icon,
			SortOrder:   category.SortOrder,
			Groups:      groupsByCategory[category.ID],
		})
	}
	return views
}

func linkView(link *model.Link, admin bool) v1.LinkView {
	view := v1.LinkView{
		ID:          link.ID,
		Title:       link.Title,
		URL:         link.URL,
		Description: link.Description,
		Tags:        link.Tags,
		Keywords:    link.Keywords,
		Kind:        link.Kind,
		Featured:    link.Featured,
		ClickCount:  link.ClickCount,
	}
	if link.LastClickedAt != nil {
		view.LastClickedAt = link.LastClickedAt.Time.UTC().Format(time.RFC3339)
	}
	if admin {
		view.CategoryID = link.CategoryID
		view.GroupID = link.GroupID
		view.Status = link.Status
		view.SortOrder = link.SortOrder
		view.HealthStatus = link.HealthStatus
		view.HealthHTTPStatus = link.HealthHTTPStatus
		view.HealthLatencyMS = link.HealthLatencyMS
		view.HealthError = link.HealthError
		if link.LastCheckedAt != nil {
			view.LastCheckedAt = link.LastCheckedAt.Time.UTC().Format(time.RFC3339)
		}
		if link.CreatedAt != nil {
			view.CreatedAt = link.CreatedAt.Time.UTC().Format(time.RFC3339)
		}
		if link.PublishedAt != nil {
			view.PublishedAt = link.PublishedAt.Time.UTC().Format(time.RFC3339)
		}
		if link.UpdatedAt != nil {
			view.UpdatedAt = link.UpdatedAt.Time.UTC().Format(time.RFC3339)
		}
	}
	return view
}
