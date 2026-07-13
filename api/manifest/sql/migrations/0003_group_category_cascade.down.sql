ALTER TABLE nav_links
  DROP CONSTRAINT fk_nav_links_group_category,
  ADD CONSTRAINT fk_nav_links_group_category
    FOREIGN KEY (group_id, category_id)
    REFERENCES nav_groups(id, category_id)
    ON DELETE RESTRICT;
