CREATE TABLE videos (
     id SERIAL primary key,
     title varchar(200),
     description varchar(4096),
     upload_date timestamp,
     userID int,
     originalSite int, -- could use an enum...
     originalLink varchar(200),
     originalID varchar(200), -- not normalized, but whatever
     newLink varchar(200)
);

CREATE TABLE tags (
    id SERIAL primary key, -- video ID
    video_id int REFERENCES videos(id),
    tag varchar(60)
);