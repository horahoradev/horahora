CREATE TABLE downloads (
    id SERIAL primary key,
    date_created timestamp,
    last_polled timestamp,
    website varchar(255),
    attribute_type varchar(255), /* tag, channel, or playlist */
    attribute_value varchar(255), /* tag id, channel id, or playlist id*/
    lock timestamp /* timestamp indicating category in use by worker, expires in 30 mins */
);

CREATE TABLE previous_downloads (
    id SERIAL primary key, -- video ID
    video_ID varchar(255),
    content_ID varchar(255), -- tag string
    upload_time timestamp,
    website varchar(255)
);

-- INSERT INTO downloads(date_created, website, attribute_type, attribute_value) VALUES (Now(), 'niconico', 'tag', 'YTPMV');
-- INSERT INTO previous_downloads(video_ID, content_ID, upload_time)