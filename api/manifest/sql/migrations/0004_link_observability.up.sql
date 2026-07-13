ALTER TABLE nav_links
  ADD COLUMN click_count BIGINT NOT NULL DEFAULT 0 CHECK (click_count >= 0),
  ADD COLUMN last_clicked_at TIMESTAMPTZ,
  ADD COLUMN health_status TEXT NOT NULL DEFAULT 'unchecked',
  ADD COLUMN last_checked_at TIMESTAMPTZ,
  ADD COLUMN health_http_status INT,
  ADD COLUMN health_latency_ms INT,
  ADD COLUMN health_error TEXT NOT NULL DEFAULT '',
  ADD CONSTRAINT ck_nav_links_health_status
    CHECK (health_status IN ('unchecked', 'healthy', 'redirected', 'broken', 'timeout', 'error'));

CREATE INDEX ix_nav_links_popular ON nav_links(status, click_count DESC, sort_order);
CREATE INDEX ix_nav_links_health ON nav_links(health_status, last_checked_at NULLS FIRST);
