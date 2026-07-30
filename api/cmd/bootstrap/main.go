// Command bootstrap installs Nav-owned initial records after schema
// migrations. It is idempotent and never overwrites operator changes.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

const initialSiteSettings = `
INSERT INTO nav_site_settings (
	id,
	name,
	title,
	description,
	search_placeholder,
	footer_tagline
)
VALUES (
	1,
	'月离导航',
	'把常用互联网，整理成工作台',
	'按任务浏览，也可以直接搜索名称、标签和域名。',
	'搜索工具、文档、社区或关键词',
	'月离导航，持续整理值得回访的互联网入口。'
)
ON CONFLICT (id) DO NOTHING`

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("NAV_DATABASE_URL"))
	if databaseURL == "" {
		fail("NAV_DATABASE_URL is required")
	}
	database, err := sql.Open("postgres", databaseURL)
	if err != nil {
		fail("open database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	if err := database.PingContext(ctx); err != nil {
		fail("connect database: %v", err)
	}
	if _, err := database.ExecContext(ctx, initialSiteSettings); err != nil {
		fail("install initial site settings: %v", err)
	}
	fmt.Println("Nav initial site records are ready")
}

func fail(format string, arguments ...any) {
	fmt.Fprintf(os.Stderr, "bootstrap: "+format+"\n", arguments...)
	os.Exit(1)
}
