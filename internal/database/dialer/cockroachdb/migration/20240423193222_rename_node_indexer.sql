-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_indexer" RENAME TO "node_worker";
ALTER TABLE "node_worker" RENAME COLUMN "worker" TO "name";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_worker" RENAME TO "node_indexer";
ALTER TABLE "node_indexer" RENAME COLUMN "name" TO "worker";
-- +goose StatementEnd
