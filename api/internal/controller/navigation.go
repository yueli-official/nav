package controller

import (
	"context"
	"time"

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
	links, err := c.service.AdminLinks(ctx, dao.LinkFilter{
		Query: req.Q, CategoryID: req.CategoryID, GroupID: req.GroupID, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	structure, err := c.service.AdminStructure(ctx)
	if err != nil {
		return nil, err
	}
	views := make([]v1.LinkView, 0, len(links))
	for _, link := range links {
		views = append(views, linkView(link, true))
	}
	return &v1.AdminListLinksRes{Links: views, Categories: categoryViews(structure, false)}, nil
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
		},
		Categories: categoryViews(value, true),
		Stats: v1.StatsView{
			CategoryCount: len(value.Categories),
			GroupCount:    len(value.Groups),
			LinkCount:     len(value.Links),
		},
	}
}

func categoryViews(value *catalog.Catalog, includeLinks bool) []v1.CategoryView {
	linksByGroup := make(map[string][]v1.LinkView, len(value.Groups))
	if includeLinks {
		for _, link := range value.Links {
			linksByGroup[link.GroupID] = append(linksByGroup[link.GroupID], linkView(link, false))
		}
	}
	groupsByCategory := make(map[string][]v1.GroupView, len(value.Categories))
	for _, group := range value.Groups {
		groupsByCategory[group.CategoryID] = append(groupsByCategory[group.CategoryID], v1.GroupView{
			ID:          group.ID,
			Title:       group.Title,
			Description: group.Description,
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
	}
	if admin {
		view.CategoryID = link.CategoryID
		view.GroupID = link.GroupID
		view.Status = link.Status
		view.SortOrder = link.SortOrder
		if link.CreatedAt != nil {
			view.CreatedAt = link.CreatedAt.Time.UTC().Format(time.RFC3339)
		}
		if link.UpdatedAt != nil {
			view.UpdatedAt = link.UpdatedAt.Time.UTC().Format(time.RFC3339)
		}
	}
	return view
}
