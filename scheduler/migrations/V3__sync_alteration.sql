ALTER TABLE downloads DROP COLUMN last_synced;

ALTER TABLE videos ADD COLUMN last_synced timestamp DEFAULT NULL;
