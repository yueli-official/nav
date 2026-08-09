ALTER TABLE nav_links
  ADD COLUMN health_check_exempt BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX ix_nav_links_health_check_exempt
  ON nav_links(health_check_exempt, health_status, last_checked_at NULLS FIRST);
