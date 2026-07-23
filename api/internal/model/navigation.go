package model

import "github.com/gogf/gf/v2/os/gtime"

type Category struct {
	ID          string `json:"id" orm:"id"`
	Title       string `json:"title" orm:"title"`
	Description string `json:"description" orm:"description"`
	Icon        string `json:"icon" orm:"icon"`
	SortOrder   int    `json:"sortOrder" orm:"sort_order"`
}

type Tag struct {
	Name      string `json:"name"`
	LinkCount int    `json:"linkCount"`
}

type SiteSettings struct {
	SearchPlaceholder string `json:"searchPlaceholder" orm:"search_placeholder"`
	RuntimeRevision   uint64 `json:"runtimeRevision" orm:"runtime_revision"`
}

type LegacySiteSettings struct {
	Name              string
	Title             string
	Description       string
	SearchPlaceholder string
	RuntimeRevision   uint64
	FooterTagline     string
}

type Group struct {
	ID          string `json:"id" orm:"id"`
	CategoryID  string `json:"categoryId" orm:"category_id"`
	Title       string `json:"title" orm:"title"`
	Description string `json:"description" orm:"description"`
	SortOrder   int    `json:"sortOrder" orm:"sort_order"`
}

type Link struct {
	ID               string      `json:"id" orm:"id"`
	CategoryID       string      `json:"categoryId" orm:"category_id"`
	GroupID          string      `json:"groupId" orm:"group_id"`
	Title            string      `json:"title" orm:"title"`
	URL              string      `json:"url" orm:"url"`
	Description      string      `json:"description" orm:"description"`
	Tags             []string    `json:"tags" orm:"tags"`
	Keywords         []string    `json:"keywords" orm:"keywords"`
	Kind             string      `json:"kind" orm:"kind"`
	Featured         bool        `json:"featured" orm:"featured"`
	Status           string      `json:"status" orm:"status"`
	SortOrder        int         `json:"sortOrder" orm:"sort_order"`
	ClickCount       int64       `json:"clickCount" orm:"click_count"`
	LastClickedAt    *gtime.Time `json:"lastClickedAt" orm:"last_clicked_at"`
	HealthStatus     string      `json:"healthStatus" orm:"health_status"`
	LastCheckedAt    *gtime.Time `json:"lastCheckedAt" orm:"last_checked_at"`
	HealthHTTPStatus int         `json:"healthHttpStatus" orm:"health_http_status"`
	HealthLatencyMS  int         `json:"healthLatencyMs" orm:"health_latency_ms"`
	HealthError      string      `json:"healthError" orm:"health_error"`
	PublishedAt      *gtime.Time `json:"publishedAt" orm:"published_at"`
	CreatedAt        *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt        *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

type LinkHealth struct {
	Status     string
	HTTPStatus int
	LatencyMS  int
	Error      string
	CheckedAt  *gtime.Time
}
