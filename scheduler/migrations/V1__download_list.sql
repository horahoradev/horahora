-- create type website_t as enum('niconico', 'bilibili', 'youtube');
create type content_type_t as enum('tag', 'channel', 'playlist');

CREATE TABLE videos (
    id SERIAL primary key, -- video ID
    video_ID varchar(255),
    url varchar(255), -- lol
    website int, -- this is denormalized, but I didn't want to have to use a trigger to ensure integrity
    upload_time timestamp
);

CREATE TABLE downloads (
    id SERIAL primary key,
    date_created timestamp,
    last_polled timestamp, /* I don't remember what this is for... FIXME? replace with order by lock desc? */
    last_synced timestamp NOT NULL DEFAULT TIMESTAMP 'epoch', /* indicates when the contents of this category were last fully queried*/
    backoff_factor int DEFAULT 1, /* minimum backoff coefficient used to determine when to fully query category */
    website int, /* this should be an enum in the future, need to write a custom scanner to read into protobuf enums */
    attribute_type content_type_t, /* tag, channel, or playlist */
    attribute_value varchar(255), /* tag id, channel id, or playlist id*/
    userID int, /* user who requested this category of content to be downloaded */
    lock timestamp, /* timestamp indicating category in use by worker, expires in 30 mins */
    UNIQUE(website, attribute_type, attribute_value)
);

CREATE TABLE downloads_to_videos (
    download_id int REFERENCES downloads(id),
    video_id int REFERENCES videos(id),
    primary key(download_id, video_id)
);

ALTER TABLE videos ADD COLUMN download_id int REFERENCES downloads(id);
ALTER TABLE videos ADD CONSTRAINT video_id_website_unique UNIQUE (video_ID, website);