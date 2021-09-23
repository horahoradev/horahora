DELETE from downloads;
DELETE from downloads_to_videos;
ALTER TABLE downloads ADD COLUMN url varchar(255) UNIQUE NOT NULL;
