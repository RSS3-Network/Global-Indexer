-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA "node";

CREATE TABLE "node"."snapshots"
(
    "date"  date   NOT NULL PRIMARY KEY,
    "count" bigint NOT NULL DEFAULT 0
);

CREATE TABLE "stake"."snapshots"
(
    "date"  date   NOT NULL PRIMARY KEY,
    "count" bigint NOT NULL DEFAULT 0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA "node";

DROP TABLE "node"."snapshots";
DROP TABLE "stake"."snapshots"
-- +goose StatementEnd
