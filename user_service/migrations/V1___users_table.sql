CREATE TABLE users (
    id SERIAL primary key, -- user ID
    username varchar(255),
    email varchar(255),
    pass_hash varchar(255),
    foreign_user_ID varchar(255),
    foreign_website varchar(255) -- Again, could use enum... lazy schema design...
);