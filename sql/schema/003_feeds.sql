-- +goose Up
CREATE TABLE IF NOT EXISTS feeds(
    id         UUID      NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_fetched_at TIMESTAMP, 
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE 
);
-- +goose Down
DROP TABLE IF EXISTS feeds;