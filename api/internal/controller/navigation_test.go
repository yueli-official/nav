package controller

import (
	"testing"

	"github.com/yueli-official/nav/api/internal/catalog"
	"github.com/yueli-official/nav/api/internal/model"
)

func TestCategoryViewsUseEmptyCollectionsInsteadOfNull(t *testing.T) {
	views := categoryViews(&catalog.Catalog{
		Categories: []*model.Category{{ID: "empty", Title: "Empty"}},
	}, false)
	if len(views) != 1 {
		t.Fatalf("views = %d, want 1", len(views))
	}
	if views[0].Groups == nil {
		t.Fatal("empty category groups must serialize as [] instead of null")
	}

	views = categoryViews(&catalog.Catalog{
		Categories: []*model.Category{{ID: "category", Title: "Category"}},
		Groups:     []*model.Group{{ID: "group", CategoryID: "category", Title: "Group"}},
	}, true)
	if len(views[0].Groups) != 1 {
		t.Fatalf("groups = %d, want 1", len(views[0].Groups))
	}
	if views[0].Groups[0].Items == nil {
		t.Fatal("empty group items must serialize as [] instead of null")
	}
}
