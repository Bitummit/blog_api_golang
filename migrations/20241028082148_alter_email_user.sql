-- +goose Up
-- +goose StatementBegin
ALTER TABLE my_user
ADD COLUMN email VARCHAR(256) DEFAULT 'ignat.001@mail.ru';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE my_user
DROP COLUMN email;
-- +goose StatementEnd