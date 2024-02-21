-- +goose Up
-- +goose StatementBegin

-- public
-- public.nodes
CREATE TABLE IF NOT EXISTS "node_info"
(
    "address"                  bytea       NOT NULL,
    "endpoint"                 text        NOT NULL,
    "is_public_good"           bool        NOT NULL,
    "stream"                   json        NOT NULL,
    "config"                   json        NOT NULL,
    "status"                   TEXT        NOT NULL DEFAULT 'offline',
    "local"                    json        NOT NULL DEFAULT '[]',
    "last_heartbeat_timestamp" timestamptz,
    "created_at"               timestamptz NOT NULL DEFAULT now(),
    "updated_at"               timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_indexes" PRIMARY KEY ("address")
);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_endpoint_unique" ON "node_info" ("endpoint");
CREATE INDEX IF NOT EXISTS "idx_is_public" ON "node_info" ("is_public_good", "created_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_created_at" ON "node_info" ("address", "created_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_status" ON "node_info" ("status");
CREATE INDEX IF NOT EXISTS "idx_last_heartbeat_timestamp" ON "node_info" ("last_heartbeat_timestamp");

-- public.checkpoints
CREATE TABLE "checkpoints"
(
    "chain_id"     bigint      NOT NULL,
    "block_number" bigint      NOT NULL,
    "block_hash"   text        NOT NULL,
    "created_at"   timestamptz NOT NULL DEFAULT now(),
    "updated_at"   timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_checkpoints" PRIMARY KEY ("chain_id")
);

-- bridge
CREATE SCHEMA "bridge";

CREATE TABLE "bridge"."transactions"
(
    "id"                text    NOT NULL,
    "type"              text    NOT NULL,
    "sender"            text    NOT NULL,
    "receiver"          text    NOT NULL,
    "token_address_l1"  text,
    "token_address_l2"  text,
    "token_value"       decimal NOT NULL,
    "data"              text,
    "chain_id"          bigint  NOT NULL,
    "block_number"      bigint,
    "transaction_index" integer,
    "block_timestamp"   timestamptz,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

CREATE INDEX "idx_transactions_sender" ON "bridge"."transactions" ("sender");
CREATE INDEX "idx_transactions_receiver" ON "bridge"."transactions" ("receiver");
CREATE INDEX "idx_transactions_address" ON "bridge"."transactions" ("sender", "receiver");
CREATE INDEX "idx_transactions_order" ON "bridge"."transactions" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);

CREATE TABLE "bridge"."events"
(
    "id"                 text        NOT NULL,
    "type"               text        NOT NULL,
    "transaction_hash"   text        NOT NULL,
    "transaction_index"  integer     NOT NULL,
    "transaction_status" integer     NOT NULL,
    "chain_id"           bigint      NOT NULL,
    "block_hash"         text        NOT NULL,
    "block_number"       bigint      NOT NULL,
    "block_timestamp"    timestamptz NOT NULL,

    CONSTRAINT "pk_events" PRIMARY KEY ("transaction_hash", "block_hash")
);

CREATE INDEX "idx_id" ON "bridge"."events" ("id");

-- stake
CREATE SCHEMA "stake";

CREATE TABLE "stake"."transactions"
(
    "id"                text        NOT NULL,
    "type"              text        NOT NULL,
    "user"              text        NOT NULL,
    "node"              text        NOT NULL,
    "value"             decimal     NOT NULL,
    "chips"             bigint[]    NOT NULL,
    "block_number"      bigint      NOT NULL,
    "transaction_index" integer     NOT NULL,
    "block_timestamp"   timestamptz NOT NULL,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

CREATE INDEX "idx_transactions_user" ON "stake"."transactions" ("user");
CREATE INDEX "idx_transactions_node" ON "stake"."transactions" ("node");
CREATE INDEX "idx_transactions_address" ON "stake"."transactions" ("user", "node");
CREATE INDEX "idx_transactions_order" ON "stake"."transactions" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);


CREATE TABLE "stake"."events"
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

CREATE INDEX "idx_id" ON "stake"."events" ("id");

CREATE TABLE "stake"."chips"
(
    "id"              decimal     NOT NULL UNIQUE,
    "owner"           text        NOT NULL,
    "node"            text        NOT NULL,
    "block_number"    bigint      NOT NULL,
    "block_timestamp" timestamptz NOT NULL,

    CONSTRAINT "pk_chips" PRIMARY KEY ("id")
);

CREATE INDEX "idx_owner" ON "stake"."chips" ("owner");
CREATE INDEX "idx_node" ON "stake"."chips" ("node");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "node_info";
DROP TABLE "checkpoints";

DROP TABLE "bridge"."transactions";
DROP TABLE "bridge"."events";
DROP SCHEMA "bridge";

DROP TABLE "stake"."transactions";
DROP TABLE "stake"."events";
DROP TABLE "stake"."chips";
DROP SCHEMA "stake";
-- +goose StatementEnd
