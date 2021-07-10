\c userservice;

INSERT INTO
  users (id, username, email, pass_hash, rank)
VALUES
  -- password: admin
  (0, 'admin', 'admin', '$2y$12$jLgDoFwdXUopJivUGqvxlurZPIdmv7I95PJ97xme35YXmeyy3gRlC', 2);
