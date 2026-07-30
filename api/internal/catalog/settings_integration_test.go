package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	_ "github.com/lib/pq"

	"platform/products/nav/api/internal/dao"
	"platform/products/nav/api/internal/navprofile"
)

func TestSiteProfileLegacyCutoverAndDowngradeProjection(t *testing.T) {
	host := os.Getenv("NAV_SITE_PROFILE_PG_HOST")
	database := os.Getenv("NAV_SITE_PROFILE_PG_DATABASE")
	if host == "" || database == "" {
		t.Skip("set NAV_SITE_PROFILE_PG_HOST and NAV_SITE_PROFILE_PG_DATABASE")
	}
	port := navTestEnvOr("NAV_SITE_PROFILE_PG_PORT", "5432")
	user := navTestEnvOr("NAV_SITE_PROFILE_PG_USER", "postgres")
	password := os.Getenv("NAV_SITE_PROFILE_PG_PASS")
	gfdb, err := gdb.New(gdb.ConfigNode{
		Type: "pgsql", Host: host, Port: port, User: user, Pass: password, Name: database,
	})
	if err != nil {
		t.Fatal(err)
	}
	dsn := (&url.URL{
		Scheme: "postgres", User: url.UserPassword(user, password),
		Host: fmt.Sprintf("%s:%s", host, port), Path: database,
		RawQuery: "sslmode=disable",
	}).String()
	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqldb.Close() })

	if _, err := sqldb.ExecContext(context.Background(), `
INSERT INTO nav_site_settings (
    id, name, title, description, search_placeholder, footer_tagline
) VALUES (1, '导航', '工作台', '站点说明', '搜索入口', '页脚说明')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    search_placeholder = EXCLUDED.search_placeholder,
    footer_tagline = EXCLUDED.footer_tagline
`); err != nil {
		t.Fatal(err)
	}
	if _, err := sqldb.ExecContext(context.Background(), `DELETE FROM site_profile_state`); err != nil {
		t.Fatal(err)
	}
	manager, err := navprofile.NewPostgres(sqldb)
	if err != nil {
		t.Fatal(err)
	}
	service := New(dao.NewPG(gfdb), Site{})
	service.SetSiteProfile(manager)
	if err := service.EnsureSiteProfile(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := service.EnsureSiteProfile(context.Background()); err != nil {
		t.Fatalf("idempotent EnsureSiteProfile: %v", err)
	}
	settings, err := service.AdminSiteSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.Snapshot.Profile.Identity.Description != "站点说明" ||
		settings.Snapshot.Profile.Footer.Tagline != "页脚说明" ||
		settings.SearchPlaceholder != "搜索入口" {
		t.Fatalf("settings = %#v", settings)
	}
	if _, err := service.SaveAdminSiteSettings(
		context.Background(), settings.Snapshot.Revision, settings.RuntimeRevision,
		settings.Snapshot.Profile, "新的搜索入口",
	); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SaveAdminSiteSettings(
		context.Background(), settings.Snapshot.Revision, settings.RuntimeRevision,
		settings.Snapshot.Profile, "过期写入",
	); err == nil {
		t.Fatal("stale Nav runtime revision did not conflict")
	}
	var name, title, description, footer any
	if err := sqldb.QueryRowContext(context.Background(), `
SELECT name, title, description, footer_tagline
FROM nav_site_settings WHERE id = 1
`).Scan(&name, &title, &description, &footer); err != nil {
		t.Fatal(err)
	}
	if name != nil || title != nil || description != nil || footer != nil {
		t.Fatalf("legacy public columns were not cleared: %#v %#v %#v %#v", name, title, description, footer)
	}

	downPath := filepath.Join("..", "..", "manifest", "sql", "migrations", "0009_site_profile_cutover.down.sql")
	downSQL, err := os.ReadFile(downPath)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := sqldb.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(context.Background(), string(downSQL)); err != nil {
		t.Fatal(err)
	}
	var restoredDescription string
	if err := tx.QueryRowContext(context.Background(), `
SELECT description FROM nav_site_settings WHERE id = 1
`).Scan(&restoredDescription); err != nil {
		t.Fatal(err)
	}
	if restoredDescription != "站点说明" {
		t.Fatalf("downgrade description = %q", restoredDescription)
	}
}

func navTestEnvOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
