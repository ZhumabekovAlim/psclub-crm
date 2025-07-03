CREATE TABLE IF NOT EXISTS pricelist_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    price_item_id INT NOT NULL,
    quantity DOUBLE NOT NULL,
    buy_price DOUBLE NOT NULL,
    total DOUBLE NOT NULL,
    user_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (price_item_id) REFERENCES price_items(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
