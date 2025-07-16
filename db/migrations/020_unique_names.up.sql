ALTER TABLE categories ADD UNIQUE KEY idx_categories_name (name);
ALTER TABLE price_items ADD UNIQUE KEY idx_price_items_name (name);
