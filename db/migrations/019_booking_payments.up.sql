CREATE TABLE IF NOT EXISTS booking_payments (
                                                id INT AUTO_INCREMENT PRIMARY KEY,
                                                booking_id INT NOT NULL,
                                                payment_type_id INT NOT NULL,
                                                amount INT NOT NULL,
                                                FOREIGN KEY (booking_id) REFERENCES bookings(id) ON DELETE CASCADE,
                                                FOREIGN KEY (payment_type_id) REFERENCES payment_types(id) ON DELETE CASCADE
);