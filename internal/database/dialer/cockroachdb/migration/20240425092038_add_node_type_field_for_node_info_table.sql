-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_info"
    ADD COLUMN "type" TEXT;

CREATE INDEX "idx_node_info_type" ON "node_info" ("type");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_info" DROP COLUMN "type";
-- +goose StatementEnd
