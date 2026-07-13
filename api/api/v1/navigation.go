package v1

import "github.com/gogf/gf/v2/frame/g"

type SiteView struct {
	Name              string `json:"name"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	SearchPlaceholder string `json:"searchPlaceholder"`
	FooterTagline     string `json:"footerTagline"`
}

type LinkView struct {
	ID          string   `json:"id"`
	CategoryID  string   `json:"categoryId,omitempty"`
	GroupID     string   `json:"groupId,omitempty"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Keywords    []string `json:"keywords,omitempty"`
	Kind        string   `json:"kind"`
	Featured    bool     `json:"featured"`
	Status      string   `json:"status,omitempty"`
	SortOrder   int      `json:"sortOrder,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
}

type GroupView struct {
	ID          string     `json:"id"`
	CategoryID  string     `json:"categoryId,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	SortOrder   int        `json:"sortOrder"`
	LinkCount   int        `json:"linkCount"`
	Items       []LinkView `json:"items"`
}

type CategoryView struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	SortOrder   int         `json:"sortOrder"`
	Groups      []GroupView `json:"groups"`
}

type StatsView struct {
	CategoryCount int `json:"categoryCount"`
	GroupCount    int `json:"groupCount"`
	LinkCount     int `json:"linkCount"`
}

type GetCatalogReq struct {
	g.Meta `path:"/api/v1/nav/catalog" method:"GET" tags:"Nav" summary:"Get the published navigation catalog"`
}

type GetCatalogRes struct {
	Version    int            `json:"version"`
	Site       SiteView       `json:"site"`
	Categories []CategoryView `json:"categories"`
	Stats      StatsView      `json:"stats"`
}

type LinkInput struct {
	CategoryID  string   `json:"categoryId" v:"required"`
	GroupID     string   `json:"groupId" v:"required"`
	Title       string   `json:"title" v:"required|length:1,200"`
	URL         string   `json:"url" v:"required|length:1,2048"`
	Description string   `json:"description" v:"required|length:1,500"`
	Tags        []string `json:"tags"`
	Keywords    []string `json:"keywords"`
	Kind        string   `json:"kind" v:"required"`
	Featured    bool     `json:"featured"`
	Status      string   `json:"status" v:"in:published,draft,archived"`
	SortOrder   int      `json:"sortOrder"`
}

type AdminListLinksReq struct {
	g.Meta     `path:"/api/v1/admin/nav/links" method:"GET" tags:"Admin Nav" summary:"List navigation links"`
	Q          string `p:"q"`
	CategoryID string `p:"categoryId"`
	GroupID    string `p:"groupId"`
	Status     string `p:"status"`
	Tag        string `p:"tag"`
	Page       int    `p:"page" d:"1" v:"min:1"`
	Size       int    `p:"size" d:"20" v:"min:1|max:100"`
}

type LifecycleCountsView struct {
	All       int `json:"all"`
	Published int `json:"published"`
	Draft     int `json:"draft"`
	Archived  int `json:"archived"`
}

type AdminListLinksRes struct {
	Links      []LinkView          `json:"links"`
	Categories []CategoryView      `json:"categories"`
	Tags       []TagView           `json:"tags"`
	Counts     LifecycleCountsView `json:"counts"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Size       int                 `json:"size"`
}

type TagView struct {
	Name      string `json:"name"`
	LinkCount int    `json:"linkCount"`
}

type AdminCreateLinkReq struct {
	g.Meta `path:"/api/v1/admin/nav/links" method:"POST" tags:"Admin Nav" summary:"Create a navigation link"`
	LinkInput
}

type AdminCreateLinkRes struct {
	Link LinkView `json:"link"`
}

type AdminUpdateLinkReq struct {
	g.Meta `path:"/api/v1/admin/nav/links/{id}" method:"PATCH" tags:"Admin Nav" summary:"Update a navigation link"`
	ID     string `p:"id" v:"required"`
	LinkInput
}

type AdminUpdateLinkRes struct {
	Link LinkView `json:"link"`
}

type AdminDeleteLinkReq struct {
	g.Meta `path:"/api/v1/admin/nav/links/{id}" method:"DELETE" tags:"Admin Nav" summary:"Delete a navigation link"`
	ID     string `p:"id" v:"required"`
}

type AdminDeleteLinkRes struct {
	Deleted bool `json:"deleted"`
}

type AdminBulkLinksReq struct {
	g.Meta `path:"/api/v1/admin/nav/links/bulk" method:"POST" tags:"Admin Nav" summary:"Bulk update navigation links"`
	IDs    []string `json:"ids" v:"required|length:1,100"`
	Action string   `json:"action" v:"required|in:publish,draft,archive,delete"`
}

type AdminBulkLinksRes struct {
	Changed   int      `json:"changed"`
	FailedIDs []string `json:"failedIds"`
}

type CategoryInput struct {
	Title       string `json:"title" v:"required|length:1,120"`
	Description string `json:"description" v:"length:0,500"`
	Icon        string `json:"icon" v:"required|length:1,120"`
	SortOrder   int    `json:"sortOrder" v:"min:0"`
}

type GroupInput struct {
	CategoryID  string `json:"categoryId" v:"required"`
	Title       string `json:"title" v:"required|length:1,120"`
	Description string `json:"description" v:"length:0,500"`
	SortOrder   int    `json:"sortOrder" v:"min:0"`
}

type AdminListStructureReq struct {
	g.Meta `path:"/api/v1/admin/nav/structure" method:"GET" tags:"Admin Nav" summary:"List categories and groups"`
}

type AdminListStructureRes struct {
	Categories []CategoryView `json:"categories"`
}

type AdminCreateCategoryReq struct {
	g.Meta `path:"/api/v1/admin/nav/categories" method:"POST" tags:"Admin Nav" summary:"Create navigation category"`
	CategoryInput
}

type AdminCreateCategoryRes struct {
	Category CategoryView `json:"category"`
}

type AdminUpdateCategoryReq struct {
	g.Meta `path:"/api/v1/admin/nav/categories/{id}" method:"PATCH" tags:"Admin Nav" summary:"Update navigation category"`
	ID     string `p:"id" v:"required"`
	CategoryInput
}

type AdminUpdateCategoryRes struct {
	Category CategoryView `json:"category"`
}

type AdminDeleteCategoryReq struct {
	g.Meta `path:"/api/v1/admin/nav/categories/{id}" method:"DELETE" tags:"Admin Nav" summary:"Delete navigation category"`
	ID     string `p:"id" v:"required"`
}

type AdminDeleteCategoryRes struct {
	Deleted bool `json:"deleted"`
}

type AdminCreateGroupReq struct {
	g.Meta `path:"/api/v1/admin/nav/groups" method:"POST" tags:"Admin Nav" summary:"Create navigation group"`
	GroupInput
}

type AdminCreateGroupRes struct {
	Group GroupView `json:"group"`
}

type AdminUpdateGroupReq struct {
	g.Meta `path:"/api/v1/admin/nav/groups/{id}" method:"PATCH" tags:"Admin Nav" summary:"Update navigation group"`
	ID     string `p:"id" v:"required"`
	GroupInput
}

type AdminUpdateGroupRes struct {
	Group GroupView `json:"group"`
}

type AdminDeleteGroupReq struct {
	g.Meta `path:"/api/v1/admin/nav/groups/{id}" method:"DELETE" tags:"Admin Nav" summary:"Delete navigation group"`
	ID     string `p:"id" v:"required"`
}

type AdminDeleteGroupRes struct {
	Deleted bool `json:"deleted"`
}

type AdminListTagsReq struct {
	g.Meta `path:"/api/v1/admin/nav/tags" method:"GET" tags:"Admin Nav" summary:"List navigation tags"`
	Q      string `p:"q"`
}

type AdminListTagsRes struct {
	Tags []TagView `json:"tags"`
}

type AdminRenameTagReq struct {
	g.Meta `path:"/api/v1/admin/nav/tags/rename" method:"POST" tags:"Admin Nav" summary:"Rename or merge a navigation tag"`
	Source string `json:"source" v:"required"`
	Target string `json:"target" v:"required"`
}

type AdminRenameTagRes struct {
	Changed int `json:"changed"`
}

type AdminDeleteTagReq struct {
	g.Meta `path:"/api/v1/admin/nav/tags/delete" method:"POST" tags:"Admin Nav" summary:"Delete a navigation tag"`
	Name   string `json:"name" v:"required"`
}

type AdminDeleteTagRes struct {
	Changed int `json:"changed"`
}

type SiteSettingsInput struct {
	Name              string `json:"name" v:"required|length:1,120"`
	Title             string `json:"title" v:"required|length:1,200"`
	Description       string `json:"description" v:"required|length:1,500"`
	SearchPlaceholder string `json:"searchPlaceholder" v:"required|length:1,200"`
	FooterTagline     string `json:"footerTagline" v:"required|length:1,300"`
}

type AdminGetSettingsReq struct {
	g.Meta `path:"/api/v1/admin/nav/settings" method:"GET" tags:"Admin Nav" summary:"Get navigation site settings"`
}

type AdminGetSettingsRes struct {
	Settings SiteView `json:"settings"`
}

type AdminUpdateSettingsReq struct {
	g.Meta `path:"/api/v1/admin/nav/settings" method:"PATCH" tags:"Admin Nav" summary:"Update navigation site settings"`
	SiteSettingsInput
}

type AdminUpdateSettingsRes struct {
	Settings SiteView `json:"settings"`
}
