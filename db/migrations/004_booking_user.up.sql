ALTER TABLE bookings ADD COLUMN user_id INT AFTER table_id;
ALTER TABLE bookings ADD CONSTRAINT fk_bookings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

