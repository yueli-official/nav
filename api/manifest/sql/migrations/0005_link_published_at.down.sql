DROP INDEX IF EXISTS ix_nav_links_published_at;

ALTER TABLE nav_links
    DROP COLUMN IF EXISTS published_at;
