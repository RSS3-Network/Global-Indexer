-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "node_invalid_response"
(
    "id"                 bigint      GENERATED BY DEFAULT AS IDENTITY (INCREMENT 1 MINVALUE 0 START 0),
    "epoch_id"           bigint      NOT NULL,
    "type"               TEXT        NOT NULL,
    "request"            TEXT        NOT NULL,
    "validator_nodes"    bytea[]     NOT NULL,
    "validator_response" json        NOT NULL,
    "node"               bytea       NOT NULL,
    "response"           json        NOT NULL,
    "created_at"         timestamptz NOT NULL DEFAULT now(),
    "updated_at"         timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pkey" PRIMARY KEY ("id")
    );

CREATE INDEX IF NOT EXISTS "idx_epoch_id" ON "node_invalid_response" ("epoch_id" DESC);
CREATE INDEX IF NOT EXISTS "idx_type" ON "node_invalid_response" ("type", "created_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_request" ON "node_invalid_response" ("request", "created_at" DESC);
CREATE INDEX IF NOT EXISTS "idx_node" ON "node_invalid_response" ("node", "created_at" DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "node_invalid_response";
-- +goose StatementEnd
