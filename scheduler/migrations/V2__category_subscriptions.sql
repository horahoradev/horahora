CREATE TABLE user_download_subscriptions (
    user_id int,
    download_id int REFERENCES downloads(id),
    primary key(user_id, download_id)
);

ALTER TABLE downloads DROP COLUMN userID;

ALTER TABLE videos ADD COLUMN dlStatus int DEFAULT 0;
/*
    0: undownloaded
    1: downloaded
    2: failed
 */