ALTER TABLE nav_site_settings
    ADD COLUMN runtime_revision BIGINT NOT NULL DEFAULT 1 CHECK (runtime_revision > 0),
    ALTER COLUMN name DROP NOT NULL,
    ALTER COLUMN title DROP NOT NULL,
    ALTER COLUMN description DROP NOT NULL,
    ALTER COLUMN footer_tagline DROP NOT NULL;

-- The application copies these fields into site_profile_state and clears them
-- in the same transaction. search_placeholder remains Nav-owned runtime data.
