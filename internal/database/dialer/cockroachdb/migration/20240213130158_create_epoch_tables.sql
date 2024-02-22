-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "epoch"
(
    "id"                      bigint      NOT NULL,
    "start_timestamp"         timestamptz NOT NULL,
    "end_timestamp"           timestamptz NOT NULL,
    "block_hash"              text        NOT NULL,
    "block_number"            bigint      NOT NULL,
    "block_timestamp"         timestamptz NOT NULL,
    "transaction_hash"        text        NOT NULL,
    "transaction_index"       bigint      NOT NULL,
    "total_operation_rewards" decimal,
    "total_staking_rewards"   decimal,
    "total_reward_items"      int,
    "created_at"              timestamptz NOT NULL DEFAULT now(),
    "updated_at"              timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_epoch" PRIMARY KEY ("transaction_hash")
);

CREATE INDEX IF NOT EXISTS "idx_timestamp" ON "epoch" ("start_timestamp" DESC, "end_timestamp" DESC);
CREATE INDEX IF NOT EXISTS "idx_epoch_id" ON "epoch" ("id" DESC, "block_number" DESC, "transaction_index" DESC);

CREATE TABLE IF NOT EXISTS "epoch_item"
(
    "epoch_id"          bigint      NOT NULL,
    "index"             int         NOT NULL,
    "node_address"      bytea       NOT NULL,
    "transaction_hash"  text        NOT NULL,
    "request_fees"      decimal     NOT NULL,
    "operation_rewards" decimal     NOT NULL,
    "staking_rewards"   decimal     NOT NULL,
    "tax_amounts"       decimal     NOT NULL,
    "created_at"        timestamptz NOT NULL DEFAULT now(),
    "updated_at"        timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_epoch_item" PRIMARY KEY ("transaction_hash" DESC, "index" ASC)
);

CREATE INDEX IF NOT EXISTS "idx_epoch_item_node_address" ON "epoch_item" ("node_address");
CREATE INDEX IF NOT EXISTS "idx_epoch_item_epoch_id" ON "epoch_item" ("epoch_id");

CREATE TABLE IF NOT EXISTS "epoch_trigger"
(
    "transaction_hash" text        NOT NULL,
    "epoch_id"         bigint      NOT NULL,
    "data"             jsonb       NOT NULL,
    "created_at"       timestamptz NOT NULL DEFAULT now(),
    "updated_at"       timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_indexes" PRIMARY KEY ("transaction_hash")
);

CREATE INDEX IF NOT EXISTS "idx_created_at" ON "epoch_trigger" ("created_at");
CREATE INDEX IF NOT EXISTS "idx_epoch_id" ON "epoch_trigger" ("epoch_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "epoch";
DROP TABLE IF EXISTS "epoch_item";
DROP TABLE IF EXISTS "epoch_trigger";
-- +goose StatementEnd
