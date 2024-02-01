-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "epoch"
(
    "id"    bigint NOT NULL,
    "start_timestamp" timestamptz NOT NULL,
    "end_timestamp"  timestamptz NOT NULL,
    "block_number" bigint NOT NULL,
    "transaction_hash" text NOT NULL,
    "transaction_index" bigint NOT NULL,
    "total_operation_rewards" decimal,
    "total_staking_rewards" decimal,
    "total_reward_items" int,
    "success" bool,
    "created_at"   timestamptz NOT NULL DEFAULT now(),
    "updated_at"   timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_epoch" PRIMARY KEY ("transaction_hash")
);

CREATE INDEX "idx_timestamp" ON "epoch" ("start_timestamp" DESC, "end_timestamp" DESC);
CREATE INDEX "idx_epoch_id" ON "epoch" ("id" DESC, "block_number" DESC, "transaction_index" DESC);

CREATE TABLE IF NOT EXISTS "epoch_item"
(
    "epoch_id" bigint NOT NULL,
    "index" int NOT NULL,
    "node_address" bytea NOT NULL,
    "request_fees" decimal NOT NULL,
    "operation_rewards" decimal NOT NULL,
    "staking_rewards" decimal NOT NULL,
    "tax_amounts" decimal NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT now(),
    "updated_at"   timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_epoch_item" PRIMARY KEY ("epoch_id" DESC, "index" ASC)
);

CREATE INDEX "idx_epoch_item_node_address" ON "epoch_item" ("node_address");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "epoch";
DROP TABLE "epoch_item";
-- +goose StatementEnd
