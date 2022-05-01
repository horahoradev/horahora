CREATE TABLE audit_events (
    id SERIAL primary key, -- user ID
    user_id int REFERENCES users(id),
    msg text,
    occurred_at timestamp DEFAULT NOW()
);
