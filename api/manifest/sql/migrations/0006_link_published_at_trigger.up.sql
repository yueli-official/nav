CREATE FUNCTION nav_set_link_published_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'published' AND NEW.published_at IS NULL THEN
        NEW.published_at = now();
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_nav_links_published_at
    BEFORE INSERT OR UPDATE OF status ON nav_links
    FOR EACH ROW
    EXECUTE FUNCTION nav_set_link_published_at();

UPDATE nav_links
SET published_at = created_at
WHERE status = 'published' AND published_at IS NULL;
