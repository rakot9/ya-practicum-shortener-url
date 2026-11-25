CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL,
    email varchar(256) NOT NULL,
    password varchar(256) NOT NULL,
    token text NULL,
    refresh_token text NULL,
    expired_at timestamp NULL,
    created_at timestamp NOT NULL
);

ALTER TABLE shortener 
ADD COLUMN user_id INTEGER,
ADD CONSTRAINT fk_url_user 
    FOREIGN KEY (user_id) 
    REFERENCES "users"(id) 
    ON DELETE CASCADE;

CREATE INDEX idx_url_user_id ON shortener(user_id);

CREATE UNIQUE INDEX idx_uq_users_uuid ON users(uuid);