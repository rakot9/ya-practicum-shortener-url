CREATE TABLE shortener (
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL,
    correlation_id VARCHAR(36) NULL,
    short_url VARCHAR(255) NOT NULL,
    original_url text NOT NULL,
    created_at timestamp NOT NULL
);

CREATE INDEX idx_uuid ON shortener(uuid);

CREATE INDEX idx_short_url ON shortener(short_url);