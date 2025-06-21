CREATE TABLE IF NOT EXISTS expense_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

ALTER TABLE expenses
    ADD COLUMN category_id INT,
    ADD COLUMN paid TINYINT(1) DEFAULT 0,
    ADD FOREIGN KEY (category_id) REFERENCES expense_categories(id) ON DELETE SET NULL;
