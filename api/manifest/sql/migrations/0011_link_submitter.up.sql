ALTER TABLE nav_links
    ADD COLUMN submitter_sub TEXT NOT NULL DEFAULT '';

CREATE INDEX ix_nav_links_submitter_scope
    ON nav_links(submitter_sub, group_id, updated_at DESC);
