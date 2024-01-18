-- +goose Up
-- +goose StatementBegin

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
DROP TABLE "bridge"."transactions";
DROP TABLE "bridge".events;
DROP SCHEMA "bridge";

DROP TABLE "stake"."stakers";
DROP SCHEMA "stake";
-- +goose StatementEnd
