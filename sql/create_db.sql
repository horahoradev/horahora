create database scheduler;
create database userservice;
create database videoservice;

\c userservice;

INSERT INTO users (username, email, pass_hash, rank) VALUES ('admin', 'admin', '$2a$05$bfl5teobsWmSk76CiG9IHuBnt94qAirXVGOnbyupVjO9K.8sOn8CK', 2)