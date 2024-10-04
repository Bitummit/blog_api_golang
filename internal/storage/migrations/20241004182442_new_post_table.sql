-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post(
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    body TEXT NOT NULL,
    author VARCHAR(128) NOT NULL
);
CREATE INDEX IF NOT EXISTS title ON post(title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post;
-- +goose StatementEnd
