-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_worker DROP CONSTRAINT "pk_indexes";

ALTER TABLE node_worker ADD CONSTRAINT "pk_node_worker" PRIMARY KEY ("epoch_id","address","network","name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE node_worker DROP CONSTRAINT "pk_node_worker";

ALTER TABLE node_worker ADD CONSTRAINT "pk_indexes" PRIMARY KEY ("address", "network", "name");
-- +goose StatementEnd
