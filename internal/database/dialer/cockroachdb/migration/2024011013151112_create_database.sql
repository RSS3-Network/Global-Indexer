-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "node_info"
(
    "address"        bytea       NOT NULL,
    "endpoint"       text        NOT NULL,
    "is_public_good" bool        NOT NULL,
    "stream"         json        NOT NULL,
    "config"         json        NOT NULL,
    "created_at"     timestamptz NOT NULL DEFAULT now(),
    "updated_at"     timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_indexes" PRIMARY KEY ("address")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_endpoint_unique" ON "node_info" ("endpoint");
CREATE INDEX IF NOT EXISTS "idx_is_public" ON "node_info" ("is_public_good", "created_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_created_at" ON "node_info" ("address", "created_at" DESC);

CREATE TABLE "checkpoints"
(
    "chain_id"     bigint      NOT NULL,
    "block_number" bigint      NOT NULL,
    "block_hash"   text        NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT now(),
    "updated_at"   timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_checkpoints" PRIMARY KEY ("chain_id")
);

CREATE SCHEMA "bridge";

CREATE TABLE "bridge"."transactions"
(
    "id"               text    NOT NULL,
    "type"             text    NOT NULL,
    "sender"           text    NOT NULL,
    receiver           text    NOT NULL,
    "token_address_l1" text,
    "token_address_l2" text,
    "token_value"      decimal NOT NULL,
    "token_decimal"    integer NOT NULL,
    "data"             text,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

CREATE TABLE "bridge"."events"
(
    "id"                 text        NOT NULL,
    "type"               text        NOT NULL,
    "transaction_hash"   text        NOT NULL,
    "transaction_index"  integer     NOT NULL,
    "transaction_status" integer     NOT NULL,
    "block_hash"         text        NOT NULL,
    "block_number"       bigint      NOT NULL,
    "block_timestamp"    timestamptz NOT NULL,

    CONSTRAINT "pk_events" PRIMARY KEY ("transaction_hash", "block_hash")
);

CREATE INDEX "idx_id" ON "bridge"."events" ("id");

CREATE SCHEMA "stake";

CREATE TABLE "stake"."stakers"
(
    "user"  text   NOT NULL,
    "node"  text   NOT NULL,
    "value" decimal NOT NULL,

    CONSTRAINT "pk_stakers" PRIMARY KEY ("user", "node")
);

CREATE INDEX "idx_value" ON "stake"."stakers" ("value");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "node_info";

DROP TABLE "bridge"."transactions";
DROP TABLE "bridge".events;
DROP SCHEMA "bridge";

DROP TABLE "stake"."stakers";
DROP SCHEMA "stake";
-- +goose StatementEnd
