-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "node_failure_response"
(
    "epoch_id"           bigint      NOT NULL,
    "status"             TEXT        NOT NULL,
    "validator_node"     bytea       NOT NULL,
    "validator_request"  TEXT        NOT NULL,
    "validator_response" json        NOT NULL,
    "verified_node"      bytea       NOT NULL,
    "verified_request"   TEXT        NOT NULL,
    "verified_response"  json        NOT NULL,
    "created_at"         timestamptz NOT NULL DEFAULT now(),
    "updated_at"         timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pkey" PRIMARY KEY ("epoch_id")

    CREATE INDEX IF NOT EXISTS "idx_validator_node" ON "node_failure_response" ("validator_node");
    CREATE INDEX IF NOT EXISTS "idx_verified_node" ON "node_failure_response" ("verified_node");
    CREATE INDEX IF NOT EXISTS "idx_status" ON "node_failure_response" ("status", "created_at" DESC);
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "node_failure_response";
-- +goose StatementEnd
