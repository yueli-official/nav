ALTER TABLE nav_links
    ADD COLUMN published_at TIMESTAMPTZ;

UPDATE nav_links
SET published_at = created_at
WHERE status = 'published' AND published_at IS NULL;

CREATE INDEX ix_nav_links_published_at
    ON nav_links (published_at DESC, title)
    WHERE status = 'published';
