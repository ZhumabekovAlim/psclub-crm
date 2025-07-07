CREATE TABLE IF NOT EXISTS repair_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

ALTER TABLE repairs
    ADD COLUMN category_id INT,
    ADD FOREIGN KEY (category_id) REFERENCES repair_categories(id) ON DELETE SET NULL;

ALTER TABLE expenses
    ADD COLUMN repair_category_id INT,
    ADD FOREIGN KEY (repair_category_id) REFERENCES repair_categories(id) ON DELETE SET NULL;
