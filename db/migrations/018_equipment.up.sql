CREATE TABLE IF NOT EXISTS equipment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    quantity DOUBLE DEFAULT 0,
    description VARCHAR(300)
);

CREATE TABLE IF NOT EXISTS equipment_inventory_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    equipment_id INT NOT NULL,
    expected DOUBLE NOT NULL,
    actual DOUBLE NOT NULL,
    difference DOUBLE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE
);
