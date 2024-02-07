-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA "node";

CREATE TABLE "node"."snapshots"
(
    "epoch_id"        bigint      NOT NULL PRIMARY KEY,
    "count"           bigint      NOT NULL DEFAULT 0,
    "block_hash"      text        NOT NULL,
    "block_number"    bigint      NOT NULL,
    "block_timestamp" timestamptz NOT NULL
);

CREATE TABLE "stake"."snapshots"
(
    "epoch_id"        bigint      NOT NULL PRIMARY KEY,
    "count"           bigint      NOT NULL DEFAULT 0,
    "block_hash"      text        NOT NULL,
    "block_number"    bigint      NOT NULL,
    "block_timestamp" timestamptz NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA "node";

DROP TABLE "node"."snapshots";
DROP TABLE "stake"."snapshots"
-- +goose StatementEnd
