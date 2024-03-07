-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "node"."events"
(
    "transaction_hash"  text        NOT NULL,
    "transaction_index" integer     NOT NULL,
    "node_id"           bigint      NOT NULL,
    "address_from"      bytea       NOT NULL,
    "address_to"        bytea       NOT NULL,
    "type"              text        NOT NULL,
    "log_index"         integer     NOT NULL,
    "chain_id"          integer     NOT NULL,
    "block_hash"        text        NOT NULL,
    "block_number"      bigint      NOT NULL,
    "block_timestamp"   timestamptz NOT NULL,
    "metadata"          jsonb       NOT NULL,
    "created_at"        timestamptz NOT NULL DEFAULT now(),
    "updated_at"        timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "events_pkey" PRIMARY KEY ("transaction_hash", "transaction_index", "log_index")
);

CREATE INDEX "events_index_node_id" ON "node"."events" ("node_id");
CREATE INDEX "events_index_address" ON "node"."events" ("address_from", "address_to");
CREATE INDEX "events_index_address_type" ON "node"."events" ("address_from", "type");
CREATE INDEX "events_index_block_number" ON "node"."events" ("block_number" DESC, "transaction_index" DESC, "log_index" DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "node"."events";
-- +goose StatementEnd
