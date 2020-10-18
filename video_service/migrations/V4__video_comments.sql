CREATE TABLE comments (
    id SERIAL primary key,
    user_id int,
    video_id int REFERENCES videos(id),
    creation_date timestamp,
    comment varchar(4096),
    parent_comment int REFERENCES comments(id)
);

CREATE TABLE comment_upvotes (
    user_id int,
    comment_id int REFERENCES comments(id),
    vote_score int,
    PRIMARY KEY(user_id, comment_id)
);