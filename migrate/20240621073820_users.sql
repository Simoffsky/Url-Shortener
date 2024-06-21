-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
