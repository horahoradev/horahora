CREATE TABLE tag_to_category (
    id SERIAL primary key,
    tag varchar(255),
    category varchar(255)
);

ALTER TABLE videos ADD COLUMN content_category varchar(255) DEFAULT "general";
