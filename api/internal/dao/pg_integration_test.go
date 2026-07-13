package dao_test

// Link DAO round-trip against a provisioned Nav PostgreSQL database. Skipped
// unless NAV_PG_HOST is set:
//
//	NAV_PG_HOST=192.168.5.5 NAV_PG_USER=postgres NAV_PG_PASS=postgres \
//	  go test -run TestPGLinkRoundTrip ./products/nav/api/internal/dao/...

import (
	"context"
	"os"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/google/uuid"

	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/model"
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
	id := "integration-" + uuid.NewString()
	link := &model.Link{
		ID: id, CategoryID: "develop", GroupID: "references",
		Title: "DAO integration " + id, URL: "https://example.com/" + id,
		Description: "Exercises GoFrame JSONB persistence", Tags: []string{"GoFrame", "JSONB"},
		Keywords: []string{"round-trip"}, Kind: "reference", Status: "draft", SortOrder: 999,
	}
	defer func() { _, _ = store.DeleteLink(ctx, id) }()

	if err := store.InsertLink(ctx, link); err != nil {
		t.Fatal(err)
	}
	items, err := store.Links(ctx, dao.LinkFilter{Query: id})
	if err != nil || len(items) != 1 || len(items[0].Tags) != 2 || items[0].Tags[1] != "JSONB" {
		t.Fatalf("insert/scan round-trip items=%#v err=%v", items, err)
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
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
