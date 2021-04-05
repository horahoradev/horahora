ALTER TABLE videos DROP COLUMN last_synced;

ALTER TABLE downloads ADD COLUMN last_synced timestamp DEFAULT NULL;
