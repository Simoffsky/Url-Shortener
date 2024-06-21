-- +goose Up

CREATE TABLE IF NOT EXISTS links (
    short_url VARCHAR(255) PRIMARY KEY,
    url TEXT NOT NULL,
    expired_at BIGINT,
    creator_login VARCHAR(255)
);



-- +goose Down
DROP TABLE IF EXISTS links;

