ALTER TABLE videos ADD COLUMN views int DEFAULT 0;
ALTER TABLE videos ADD COLUMN rating float DEFAULT 0.00;

CREATE TABLE ratings (
    user_id int,
    video_id int REFERENCES videos(id),
    rating float,
    PRIMARY KEY(user_id, video_id)
);