package dao_test

// Link DAO round-trip against a provisioned Nav PostgreSQL database. Skipped
// unless NAV_PG_HOST is set:
//
//	NAV_PG_HOST=192.168.5.5 NAV_PG_USER=postgres NAV_PG_PASS=postgres \
//	  go test -run TestPGLinkRoundTrip ./products/nav/api/internal/dao/...

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/yueli-official/foundation/go/identifier"
	"github.com/yueli-official/nav/api/internal/dao"
	"github.com/yueli-official/nav/api/internal/model"
)

func TestPGLinkRoundTrip(t *testing.T) {
	host := os.Getenv("NAV_PG_HOST")
	if host == "" {
		t.Skip("set NAV_PG_HOST to run the Nav DAO integration test")
	}
	db, err := gdb.New(gdb.ConfigNode{
		Type: "pgsql",
		Host: host,
		Port: envOr("NAV_PG_PORT", "5432"),
		User: envOr("NAV_PG_USER", "postgres"),
		Pass: os.Getenv("NAV_PG_PASS"),
		Name: envOr("NAV_PG_DATABASE", "nav_yueli"),
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	store := dao.NewPG(db)
	id := "integration-" + identifier.MustNew().String()
	tagOriginal := "integration-tag-" + identifier.MustNew().String()
	tagRenamed := tagOriginal + "-renamed"
	link := &model.Link{
		ID: id, CategoryID: "develop", GroupID: "references",
		Title: "DAO integration " + id, URL: "https://example.com/" + id,
		Description: "Exercises GoFrame JSONB persistence", Tags: []string{tagOriginal, "JSONB"},
		Keywords: []string{"round-trip"}, Kind: "reference", Status: "draft", SortOrder: 999,
	}
	defer func() { _, _ = store.DeleteLink(ctx, id) }()

	if err := store.InsertLink(ctx, link); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if err := store.UpsertFavicon(ctx, &model.FaviconCache{
		LinkID: id, SourceURL: link.URL, Content: []byte("integration-png"),
		ContentType: "image/png", ContentHash: strings.Repeat("a", 64),
		FetchedAt: gtime.NewFromTime(now), RefreshAfter: gtime.NewFromTime(now.Add(24 * time.Hour)),
		LastAttemptAt: gtime.NewFromTime(now),
	}); err != nil {
		t.Fatal(err)
	}
	items, err := store.Links(ctx, dao.LinkFilter{Query: id})
	if err != nil || len(items) != 1 || len(items[0].Tags) != 2 || items[0].Tags[1] != "JSONB" || items[0].FaviconRevision != strings.Repeat("a", 64) {
		t.Fatalf("insert/scan round-trip items=%#v err=%v", items, err)
	}
	exempted, err := store.UpdateLinkCheckExempt(ctx, id, true)
	if err != nil || !exempted {
		t.Fatalf("set health check exemption changed=%v err=%v", exempted, err)
	}
	exemptItems, err := store.Links(ctx, dao.LinkFilter{Health: "exempt", Query: id})
	if err != nil || len(exemptItems) != 1 || !exemptItems[0].HealthCheckExempt {
		t.Fatalf("exempt filter items=%#v err=%v", exemptItems, err)
	}
	checkableCount, err := store.CountLinks(ctx, dao.LinkFilter{Query: id, ExcludeHealthCheckExempt: true})
	if err != nil || checkableCount != 0 {
		t.Fatalf("checkable count=%d err=%v, want 0", checkableCount, err)
	}
	if _, err := store.UpdateLinkCheckExempt(ctx, id, false); err != nil {
		t.Fatal(err)
	}
	cached, err := store.FaviconByLinkID(ctx, id)
	if err != nil || cached == nil || string(cached.Content) != "integration-png" {
		t.Fatalf("favicon cache=%#v err=%v", cached, err)
	}
	count, err := store.CountLinks(ctx, dao.LinkFilter{Tag: tagOriginal})
	if err != nil || count != 1 {
		t.Fatalf("tag count=%d err=%v", count, err)
	}
	tags, err := store.Tags(ctx, tagOriginal)
	if err != nil || len(tags) != 1 || tags[0].Name != tagOriginal {
		t.Fatalf("tag list=%#v err=%v", tags, err)
	}
	changed, err := store.RenameTag(ctx, tagOriginal, tagRenamed)
	if err != nil || changed != 1 {
		t.Fatalf("rename tag changed=%d err=%v", changed, err)
	}
	changed, err = store.DeleteTag(ctx, tagRenamed)
	if err != nil || changed != 1 {
		t.Fatalf("delete tag changed=%d err=%v", changed, err)
	}
	changed, err = store.BulkUpdateLinks(ctx, []string{id}, "published")
	if err != nil || changed != 1 {
		t.Fatalf("bulk update changed=%d err=%v", changed, err)
	}

	link.Title = "DAO updated " + id
	link.Tags = []string{"updated"}
	updated, err := store.UpdateLink(ctx, link)
	if err != nil || !updated {
		t.Fatalf("update result=%v err=%v", updated, err)
	}
	items, err = store.Links(ctx, dao.LinkFilter{Query: id})
	if err != nil || len(items) != 1 || items[0].Title != link.Title || len(items[0].Tags) != 1 {
		t.Fatalf("update/scan round-trip items=%#v err=%v", items, err)
	}

	deleted, err := store.DeleteLink(ctx, id)
	if err != nil || !deleted {
		t.Fatalf("delete result=%v err=%v", deleted, err)
	}
	exists, err := store.LinkExists(ctx, id)
	if err != nil || exists {
		t.Fatalf("link still exists=%v err=%v", exists, err)
	}
	cached, err = store.FaviconByLinkID(ctx, id)
	if err != nil || cached != nil {
		t.Fatalf("favicon cache did not cascade with link delete: cache=%#v err=%v", cached, err)
	}

	category := &model.Category{ID: "integration-category-" + identifier.MustNew().String(), Title: "Integration category", Icon: "i-tabler-folder", SortOrder: 999}
	categoryTarget := &model.Category{ID: "integration-category-" + identifier.MustNew().String(), Title: "Integration target category", Icon: "i-tabler-folder", SortOrder: 999}
	group := &model.Group{ID: "integration-group-" + identifier.MustNew().String(), CategoryID: category.ID, Title: "Integration group", SortOrder: 999}
	groupLink := &model.Link{ID: "integration-group-link-" + identifier.MustNew().String(), CategoryID: category.ID, GroupID: group.ID, Title: "Group move fixture", URL: "https://example.com/group-move", Description: "Verifies deferred topic moves", Kind: "reference", Status: "draft"}
	defer func() {
		_, _ = store.DeleteLink(ctx, groupLink.ID)
		_, _ = store.DeleteGroup(ctx, group.ID)
		_, _ = store.DeleteCategory(ctx, category.ID)
		_, _ = store.DeleteCategory(ctx, categoryTarget.ID)
	}()
	if err := store.InsertCategory(ctx, category); err != nil {
		t.Fatal(err)
	}
	if err := store.InsertCategory(ctx, categoryTarget); err != nil {
		t.Fatal(err)
	}
	if err := store.InsertGroup(ctx, group); err != nil {
		t.Fatal(err)
	}
	if err := store.InsertLink(ctx, groupLink); err != nil {
		t.Fatal(err)
	}
	category.Title = "Updated integration category"
	if updated, updateErr := store.UpdateCategory(ctx, category); updateErr != nil || !updated {
		t.Fatalf("category update=%v err=%v", updated, updateErr)
	}
	group.Title = "Updated integration group"
	group.CategoryID = categoryTarget.ID
	if updated, updateErr := store.UpdateGroup(ctx, group); updateErr != nil || !updated {
		t.Fatalf("group update=%v err=%v", updated, updateErr)
	}
	movedLinks, moveErr := store.Links(ctx, dao.LinkFilter{GroupID: group.ID})
	if moveErr != nil || len(movedLinks) != 1 || movedLinks[0].CategoryID != categoryTarget.ID {
		t.Fatalf("group move links=%#v err=%v", movedLinks, moveErr)
	}
	if deleted, deleteErr := store.DeleteLink(ctx, groupLink.ID); deleteErr != nil || !deleted {
		t.Fatalf("group fixture link delete=%v err=%v", deleted, deleteErr)
	}
	if deleted, deleteErr := store.DeleteGroup(ctx, group.ID); deleteErr != nil || !deleted {
		t.Fatalf("group delete=%v err=%v", deleted, deleteErr)
	}
	if deleted, deleteErr := store.DeleteCategory(ctx, category.ID); deleteErr != nil || !deleted {
		t.Fatalf("category delete=%v err=%v", deleted, deleteErr)
	}
	if deleted, deleteErr := store.DeleteCategory(ctx, categoryTarget.ID); deleteErr != nil || !deleted {
		t.Fatalf("target category delete=%v err=%v", deleted, deleteErr)
	}

	settings, err := store.SiteSettings(ctx)
	if err != nil || settings == nil || settings.SearchPlaceholder == "" {
		t.Fatalf("site settings=%#v err=%v", settings, err)
	}
	if err := store.UpsertSiteSettings(ctx, settings); err != nil {
		t.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
