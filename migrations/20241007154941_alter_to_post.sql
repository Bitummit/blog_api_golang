-- +goose Up
-- +goose StatementBegin
ALTER TABLE post
DROP COLUMN author;
ALTER TABLE post
ADD COLUMN author_id INT;
ALTER TABLE post
ADD CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES author(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE post
DROP COLUMN author_id;
-- +goose StatementEnd
