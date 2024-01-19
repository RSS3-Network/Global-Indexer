-- +goose Up
-- +goose StatementBegin
CREATE TABLE "stake"."transactions"
(
    "id"    text      NOT NULL,
    "type"  text      NOT NULL,
    "user"  text      NOT NULL,
    "node"  text      NOT NULL,
    "value" decimal   NOT NULL,
    "chips" bigint[] NOT NULL,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stake"."transactions";
DROP TABLE "stake".events;
-- +goose StatementEnd
