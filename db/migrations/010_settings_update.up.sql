ALTER TABLE settings
    ADD COLUMN tables_count INT DEFAULT 0,
    ADD COLUMN notification_time INT DEFAULT 0;
