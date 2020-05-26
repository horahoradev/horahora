CREATE TABLE videos (
     id SERIAL primary key,
     title varchar(60),
     description varchar(4096),
     userID int,
     originalSite int, -- could use an enum...
     originalLink varchar(60),
     newLink varchar(60)
);

CREATE TABLE tags (
    id SERIAL primary key, -- video ID
    video_id int REFERENCES videos(id),
    tag varchar(60)
);