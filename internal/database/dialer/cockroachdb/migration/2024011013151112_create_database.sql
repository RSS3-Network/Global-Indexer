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

CREATE TABLE IF NOT EXISTS "node_stat"
(
    "address"                     bytea       NOT NULL,
    "endpoint"                    text        NOT NULL,
    "points"                      decimal     NOT NULL,
    "is_public_good"              bool        NOT NULL,
    "is_full_node"                bool        NOT NULL,
    "is_rss_node"                 bool        NOT NULL,
    "staking"                     decimal     NOT NULL,
    "total_request_count"         int         NOT NULL,
    "epoch_request_count"         int         NOT NULL,
    "epoch_invalid_request_count" int         NOT NULL,
    "decentralized_network_count" int         NOT NULL,
    "federated_network_count"     int         NOT NULL,
    "indexer_count"               int         NOT NULL,
    "reset_at"                    timestamptz NOT NULL,
    "created_at"                  timestamptz NOT NULL DEFAULT now(),
    "updated_at"                  timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pk_indexes" PRIMARY KEY ("address")
);

CREATE INDEX IF NOT EXISTS "idx_indexes_points" ON "node_stat" ("points" DESC);
CREATE INDEX IF NOT EXISTS "idx_indexes_is_full_node" ON "node_stat" ("is_full_node", "points" DESC);
CREATE INDEX IF NOT EXISTS "idx_indexes_is_rss_node" ON "node_stat" ("is_rss_node", "points" DESC);

CREATE TABLE IF NOT EXISTS "node_indexer"
(
    "address" bytea NOT NULL,
    "network" text  NOT NULL,
    "worker"  text  NOT NULL,

    CONSTRAINT "pk_indexes" PRIMARY KEY ("address","network","worker")
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "node_info";
DROP TABLE IF EXISTS "node_stat";
DROP TABLE IF EXISTS "node_indexer";
-- +goose StatementEnd
