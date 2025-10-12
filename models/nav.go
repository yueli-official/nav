package models

type Collection struct {
	Title       string   `yaml:"title" json:"title"`
	Link        string   `yaml:"link" json:"link"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// Group 分组
type Group struct {
	GID   string `yaml:"gid" json:"gid"`
	Title string `yaml:"title" json:"title"`
	Order int    `yaml:"order" json:"order"`
}

// GroupWithCollections 带收藏项的分组
type GroupWithCollections struct {
	Group       Group        `yaml:"group" json:"group"`
	Collections []Collection `yaml:"collections" json:"collections"`
}

// Category 分类
type Category struct {
	CID   string `yaml:"cid" json:"cid"`
	Title string `yaml:"title" json:"title"`
	Order int    `yaml:"order" json:"order"`
}

// CategoryData 完整的分类数据
type CategoryData struct {
	Category Category               `yaml:"category" json:"category"`
	Groups   []GroupWithCollections `yaml:"groups" json:"groups"`
}
