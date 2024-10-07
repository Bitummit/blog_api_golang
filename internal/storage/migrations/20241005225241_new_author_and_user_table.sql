-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS my_user(
    id SERIAL PRIMARY KEY,
    username VARCHAR(256) NOT NULL,
    pass VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS author(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES my_user(id),
    first_name VARCHAR(256) NOT NULL,
    last_name VARCHAR(256) NOT NULL,
    age INT
);

CREATE TABLE IF NOT EXISTS token (
    id SERIAL PRIMARY KEY,
    access_token VARCHAR(512) NOT NULL,
    refresh_token VARCHAR(512) NOT NULL
);

CREATE INDEX IF NOT EXISTS last_name ON author(last_name);
CREATE INDEX IF NOT EXISTS access_token ON token(access_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE author;
DROP TABLE token;
-- +goose StatementEnd
