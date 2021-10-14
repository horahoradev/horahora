CREATE TABLE archival_events (
    id SERIAL primary key,
    video_url varchar(255),
    parent_url varchar(255),
    download_id int REFERENCES downloads(id),
    event_message varchar(255),
    event_time timestamp
);