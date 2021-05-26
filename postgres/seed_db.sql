\c userservice;

INSERT INTO
  users (username, email, pass_hash, rank)
VALUES
  -- password: admin
  ('admin', 'admin', '$2y$12$jLgDoFwdXUopJivUGqvxlurZPIdmv7I95PJ97xme35YXmeyy3gRlC', 2);
