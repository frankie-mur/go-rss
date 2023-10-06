-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows(
    id         UUID      NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);
-- +goose Down
DROP TABLE IF EXISTS feed_follows;