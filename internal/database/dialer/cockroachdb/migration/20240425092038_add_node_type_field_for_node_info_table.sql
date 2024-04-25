-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_info"
    ADD COLUMN "type" TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_info" DROP COLUMN "type";
-- +goose StatementEnd
