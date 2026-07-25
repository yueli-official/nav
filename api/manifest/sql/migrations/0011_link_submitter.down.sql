DROP INDEX IF EXISTS ix_nav_links_submitter_scope;
ALTER TABLE nav_links DROP COLUMN IF EXISTS submitter_sub;
