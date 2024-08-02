-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_info"
    ADD COLUMN "access_token" TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_info"
    DROP COLUMN "access_token";
-- +goose StatementEnd
