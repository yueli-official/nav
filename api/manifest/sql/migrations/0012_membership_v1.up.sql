CREATE TABLE nav_memberships (
    id                      UUID PRIMARY KEY,
    user_key                TEXT NOT NULL UNIQUE,
    status                  TEXT NOT NULL DEFAULT 'active',
    display_name_snapshot   TEXT NOT NULL DEFAULT '',
    handle_snapshot         TEXT NOT NULL DEFAULT '',
    avatar_media_key_snapshot TEXT NOT NULL DEFAULT '',
    profile_synced_at       TIMESTAMPTZ,
    join_reconciled_at      TIMESTAMPTZ,
    joined_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    suspended_at            TIMESTAMPTZ,
    suspended_by_user_key   TEXT NOT NULL DEFAULT '',
    suspension_reason       TEXT NOT NULL DEFAULT '',
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_nav_memberships_user_key
        CHECK (user_key ~ '^[1-9A-HJ-NP-Za-km-z]{8}$'),
    CONSTRAINT ck_nav_memberships_status
        CHECK (status IN ('active', 'suspended'))
);

CREATE INDEX ix_nav_memberships_status_activity
    ON nav_memberships(status, last_seen_at DESC, id);
CREATE INDEX ix_nav_memberships_display_name
    ON nav_memberships(lower(display_name_snapshot));
CREATE INDEX ix_nav_memberships_handle
    ON nav_memberships(lower(handle_snapshot));

CREATE TABLE nav_membership_events (
    id              UUID PRIMARY KEY,
    membership_id   UUID NOT NULL REFERENCES nav_memberships(id) ON DELETE RESTRICT,
    actor_user_key  TEXT NOT NULL DEFAULT '',
    event_type      TEXT NOT NULL,
    previous_status TEXT NOT NULL DEFAULT '',
    next_status     TEXT NOT NULL DEFAULT '',
    reason          TEXT NOT NULL DEFAULT '',
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_nav_membership_events_type
        CHECK (event_type IN ('joined', 'status_changed')),
    CONSTRAINT ck_nav_membership_events_previous_status
        CHECK (previous_status IN ('', 'active', 'suspended')),
    CONSTRAINT ck_nav_membership_events_next_status
        CHECK (next_status IN ('active', 'suspended'))
);

CREATE INDEX ix_nav_membership_events_membership_time
    ON nav_membership_events(membership_id, occurred_at DESC, id DESC);
