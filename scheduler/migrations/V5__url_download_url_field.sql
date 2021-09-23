DELETE from downloads_to_videos;
DELETE from user_download_subscriptions;
DELETE from downloads;
ALTER TABLE downloads ADD COLUMN url varchar(255) UNIQUE NOT NULL;
