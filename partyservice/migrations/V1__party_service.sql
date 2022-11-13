CREATE TABLE parties (
    id int primary key,
    LeaderID int
);

CREATE TABLE watchers (
    PartyID int REFERENCES parties(id),
    UserID int, -- no foreign key here
    Username varchar(255),
    primary key(PartyID, UserID);
);

CREATE TABLE video_queue (
    id SERIAL primary key,
    PartyID int REFERENCERS parties(id),
    TS timestamp
    VideoID int,
    Title text,
    Location text
);
