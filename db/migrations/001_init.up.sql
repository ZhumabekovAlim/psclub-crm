-- 1. Клиенты
CREATE TABLE IF NOT EXISTS clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    date_of_birth DATE,
    channel VARCHAR(50),
    bonus INT DEFAULT 0,
    visits INT DEFAULT 0,
    income INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 2. Сотрудники (Users/Admins)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    salary_hookah DOUBLE DEFAULT 0,
    salary_bar DOUBLE DEFAULT 0,
    salary_shift INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 3. Категории и подкатегории товаров/услуг
CREATE TABLE IF NOT EXISTS categories (
                                          id INT AUTO_INCREMENT PRIMARY KEY,
                                          name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS subcategories (
                                             id INT AUTO_INCREMENT PRIMARY KEY,
                                             category_id INT NOT NULL,
                                             name VARCHAR(50) NOT NULL,
                                             FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- 4. Прайс-лист (товары/услуги)
CREATE TABLE IF NOT EXISTS price_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category_id INT NOT NULL,
    subcategory_id INT,
    quantity DOUBLE DEFAULT 0,
    sale_price DOUBLE NOT NULL,
    buy_price DOUBLE DEFAULT 0,
    is_set TINYINT(1) DEFAULT 0,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (subcategory_id) REFERENCES subcategories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS price_sets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category_id INT NOT NULL,
    subcategory_id INT,
    price INT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (subcategory_id) REFERENCES subcategories(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS set_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    price_set_id INT NOT NULL,
    item_id INT NOT NULL,
    quantity DOUBLE NOT NULL,
    FOREIGN KEY (price_set_id) REFERENCES price_sets(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES price_items(id) ON DELETE CASCADE
);

-- 5. Категории и столы
CREATE TABLE IF NOT EXISTS table_categories (
                                                id INT AUTO_INCREMENT PRIMARY KEY,
                                                name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS tables (
                                      id INT AUTO_INCREMENT PRIMARY KEY,
                                      category_id INT NOT NULL,
                                      name VARCHAR(50) NOT NULL,
                                      FOREIGN KEY (category_id) REFERENCES table_categories(id) ON DELETE CASCADE
);

-- 6. Бронирования и позиции бронирования
CREATE TABLE IF NOT EXISTS bookings (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        client_id INT NOT NULL,
                                        table_id INT NOT NULL,
                                        start_time DATETIME NOT NULL,
                                        end_time DATETIME NOT NULL,
                                        note VARCHAR(300),
                                        discount INT DEFAULT 0,
                                        discount_reason VARCHAR(300),
                                        total_amount INT NOT NULL,
                                        bonus_used INT DEFAULT 0,
                                        payment_status VARCHAR(30) NOT NULL,
                                        payment_type_id INT NOT NULL,
                                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                        FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
                                        FOREIGN KEY (table_id) REFERENCES tables(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS booking_items (
                                             id INT AUTO_INCREMENT PRIMARY KEY,
                                             booking_id INT NOT NULL,
                                             item_id INT NOT NULL,
                                             quantity DOUBLE NOT NULL,
                                             price INT NOT NULL,
                                             discount INT DEFAULT 0,
                                             created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                             FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
                                             FOREIGN KEY (item_id) REFERENCES price_items(id) ON DELETE CASCADE
);

-- 7. Расходы
CREATE TABLE IF NOT EXISTS expenses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATETIME NOT NULL,
    title VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    total DOUBLE NOT NULL,
    description VARCHAR(300),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 8. Ремонты
CREATE TABLE IF NOT EXISTS repairs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    date DATETIME NOT NULL,
    color VARCHAR(30),
    vin VARCHAR(50),
    description VARCHAR(300),
    price DOUBLE DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 9. История закупа/прихода на склад
CREATE TABLE IF NOT EXISTS price_item_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    price_item_id INT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    quantity DOUBLE NOT NULL,
    buy_price DOUBLE DEFAULT 0,
    total DOUBLE NOT NULL,
    user_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (price_item_id) REFERENCES price_items(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 10. Касса
CREATE TABLE IF NOT EXISTS cashbox (
    id INT AUTO_INCREMENT PRIMARY KEY,
    amount DOUBLE NOT NULL DEFAULT 0
);
INSERT IGNORE INTO cashbox (id, amount) VALUES (1, 0);

-- 11. Глобальные настройки
CREATE TABLE IF NOT EXISTS settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    payment_type VARCHAR(30),
    block_time INT DEFAULT 0,
    bonus_percent INT DEFAULT 0,
    work_time_from TIME,
    work_time_to TIME
);

-- 12. Типы оплат
CREATE TABLE IF NOT EXISTS payment_types (
                                             id INT AUTO_INCREMENT PRIMARY KEY,
                                             name VARCHAR(30) NOT NULL
);
