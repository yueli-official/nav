DROP INDEX IF EXISTS ix_nav_links_health_check_exempt;

ALTER TABLE nav_links
  DROP COLUMN IF EXISTS health_check_exempt;
