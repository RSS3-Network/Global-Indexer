-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_worker ADD COLUMN "epoch_id" bigint NOT NULL;

ALTER TABLE node_worker ADD COLUMN "is_active" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS "idx_node_worker_is_active" ON node_worker(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE node_worker DROP COLUMN "epoch_id";

ALTER TABLE node_worker DROP COLUMN "is_active";

DROP INDEX "idx_node_worker_is_active";
-- +goose StatementEnd
