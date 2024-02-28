-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_info"
    ADD "avatar" jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_info"
    DROP COLUMN "avatar";
-- +goose StatementEnd
