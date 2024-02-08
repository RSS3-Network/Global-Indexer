-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "node_info"
(
    "address"                  bytea       NOT NULL,
    "endpoint"                 text        NOT NULL,
    "is_public_good"           bool        NOT NULL,
    "stream"                   json        NOT NULL,
    "config"                   json        NOT NULL,
    "status"                   TEXT        NOT NULL DEFAULT 'offline',
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "node_info";
-- +goose StatementEnd
