-- create type website_t as enum('niconico', 'bilibili', 'youtube');
create type content_type_t as enum('tag', 'channel', 'playlist');


CREATE TABLE downloads (
    id SERIAL primary key,
    date_created timestamp,
    last_polled timestamp,
    website int, /* this should be an enum in the future, need to write a custom scanner to read into protobuf enums */
    attribute_type content_type_t, /* tag, channel, or playlist */
    attribute_value varchar(255), /* tag id, channel id, or playlist id*/
    userID int, /* user who requested this category of content to be downloaded */
    lock timestamp, /* timestamp indicating category in use by worker, expires in 30 mins */
    UNIQUE(website, attribute_type, attribute_value)
);

CREATE TABLE previous_downloads (
    id SERIAL primary key, -- video ID
    video_ID varchar(255),
    content_ID varchar(255), -- tag string
    upload_time timestamp,
    website varchar(255)
);
