CREATE TABLE nav_link_favicons (
    link_id          TEXT PRIMARY KEY REFERENCES nav_links(id) ON DELETE CASCADE,
    source_url       TEXT NOT NULL,
    content          BYTEA,
    content_type     TEXT NOT NULL DEFAULT '',
    content_hash     TEXT NOT NULL DEFAULT '',
    fetched_at       TIMESTAMPTZ,
    refresh_after    TIMESTAMPTZ NOT NULL,
    last_attempt_at  TIMESTAMPTZ NOT NULL,
    last_error       TEXT NOT NULL DEFAULT '',
    CONSTRAINT ck_nav_link_favicons_payload
        CHECK (
            (content IS NULL AND content_type = '' AND content_hash = '' AND fetched_at IS NULL)
            OR
            (content IS NOT NULL AND octet_length(content) > 0 AND content_type <> ''
                AND content_hash ~ '^[0-9a-f]{64}$' AND fetched_at IS NOT NULL)
        )
);

CREATE INDEX ix_nav_link_favicons_refresh
    ON nav_link_favicons(refresh_after, link_id);
