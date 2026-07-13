DROP INDEX IF EXISTS ix_nav_links_health;
DROP INDEX IF EXISTS ix_nav_links_popular;

ALTER TABLE nav_links
  DROP CONSTRAINT IF EXISTS ck_nav_links_health_status,
  DROP COLUMN IF EXISTS health_error,
  DROP COLUMN IF EXISTS health_latency_ms,
  DROP COLUMN IF EXISTS health_http_status,
  DROP COLUMN IF EXISTS last_checked_at,
  DROP COLUMN IF EXISTS health_status,
  DROP COLUMN IF EXISTS last_clicked_at,
  DROP COLUMN IF EXISTS click_count;
