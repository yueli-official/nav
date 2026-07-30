UPDATE nav_site_settings
SET
    name = profile.document #>> '{identity,name}',
    title = COALESCE(profile.document #>> '{identity,tagline}', ''),
    description = COALESCE(profile.document #>> '{identity,description}', ''),
    footer_tagline = COALESCE(profile.document #>> '{footer,tagline}', '')
FROM site_profile_state AS profile
WHERE nav_site_settings.id = 1
  AND profile.id = 1;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM nav_site_settings
        WHERE name IS NULL
           OR title IS NULL
           OR description IS NULL
           OR footer_tagline IS NULL
    ) THEN
        RAISE EXCEPTION 'cannot downgrade Nav Site Profile without a complete profile projection';
    END IF;
END $$;

ALTER TABLE nav_site_settings
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN title SET NOT NULL,
    ALTER COLUMN description SET NOT NULL,
    ALTER COLUMN footer_tagline SET NOT NULL;

ALTER TABLE nav_site_settings DROP COLUMN runtime_revision;
