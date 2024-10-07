-- +goose Up
-- +goose StatementBegin
ALTER TABLE my_user
ADD CONSTRAINT unique_username UNIQUE(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
