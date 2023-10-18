-- +goose Up
ALTER TABLE users
    ADD email TEXT UNIQUE;

-- +goose Down
ALTER TABLE users
    DROP COLUMN email;
