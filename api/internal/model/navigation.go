package model

import "github.com/gogf/gf/v2/os/gtime"

type Category struct {
	ID          string `json:"id" orm:"id"`
	Title       string `json:"title" orm:"title"`
	Description string `json:"description" orm:"description"`
	Icon        string `json:"icon" orm:"icon"`
	SortOrder   int    `json:"sortOrder" orm:"sort_order"`
}

type Group struct {
	ID          string `json:"id" orm:"id"`
	CategoryID  string `json:"categoryId" orm:"category_id"`
	Title       string `json:"title" orm:"title"`
	Description string `json:"description" orm:"description"`
	SortOrder   int    `json:"sortOrder" orm:"sort_order"`
}

type Link struct {
	ID          string      `json:"id" orm:"id"`
	CategoryID  string      `json:"categoryId" orm:"category_id"`
	GroupID     string      `json:"groupId" orm:"group_id"`
	Title       string      `json:"title" orm:"title"`
	URL         string      `json:"url" orm:"url"`
	Description string      `json:"description" orm:"description"`
	Tags        []string    `json:"tags" orm:"tags"`
	Keywords    []string    `json:"keywords" orm:"keywords"`
	Kind        string      `json:"kind" orm:"kind"`
	Featured    bool        `json:"featured" orm:"featured"`
	Status      string      `json:"status" orm:"status"`
	SortOrder   int         `json:"sortOrder" orm:"sort_order"`
	CreatedAt   *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt" orm:"updated_at"`
}
