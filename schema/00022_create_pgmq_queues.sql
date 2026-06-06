-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgmq CASCADE;
SELECT pgmq.create('notifications');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT pgmq.drop('notifications');
-- +goose StatementEnd
