CREATE TABLE IF NOT EXISTS inventory_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    price_item_id INT NOT NULL,
    expected DOUBLE NOT NULL,
    actual DOUBLE NOT NULL,
    difference DOUBLE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (price_item_id) REFERENCES price_items(id) ON DELETE CASCADE
);
