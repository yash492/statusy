-- +goose Up
-- +goose StatementBegin
SELECT pgmq.create('notifications');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT pgmq.drop('notifications');
-- +goose StatementEnd
