ALTER TABLE videos ADD COLUMN is_approved bool DEFAULT FALSE;

CREATE TABLE approvals (
    user_id int,
    video_id int REFERENCES videos(id),
    PRIMARY KEY(user_id, video_id)
);