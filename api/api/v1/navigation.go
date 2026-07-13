package v1

import "github.com/gogf/gf/v2/frame/g"

type SiteView struct {
	Name              string `json:"name"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	SearchPlaceholder string `json:"searchPlaceholder"`
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
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Items       []LinkView `json:"items"`
}

type CategoryView struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
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
}

type AdminListLinksRes struct {
	Links      []LinkView     `json:"links"`
	Categories []CategoryView `json:"categories"`
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
