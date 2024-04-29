-- +goose Up
-- create "average_tax_rate_submissions" table
CREATE TABLE "public"."average_tax_rate_submissions" (
  "epoch_id" bigint NOT NULL,
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "average_tax_rate" numeric NOT NULL,
  "transaction_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("epoch_id")
);
-- create index "average_tax_rate_submissions_id_idx" to table: "average_tax_rate_submissions"
CREATE INDEX "average_tax_rate_submissions_id_idx" ON "public"."average_tax_rate_submissions" ("id" DESC);
-- create index "average_tax_rate_submissions_transaction_hash_idx" to table: "average_tax_rate_submissions"
CREATE INDEX "average_tax_rate_submissions_transaction_hash_idx" ON "public"."average_tax_rate_submissions" ("transaction_hash");
-- create "bridge_events" table
CREATE TABLE "public"."bridge_events" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NOT NULL,
  "transaction_status" bigint NOT NULL,
  "chain_id" bigint NOT NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  PRIMARY KEY ("transaction_hash", "block_hash")
);
-- create index "idx_id" to table: "bridge_events"
CREATE INDEX "idx_id" ON "public"."bridge_events" ("id");
-- create "bridge_transactions" table
CREATE TABLE "public"."bridge_transactions" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "sender" text NOT NULL,
  "receiver" text NOT NULL,
  "token_address_l1" text NULL,
  "token_address_l2" text NULL,
  "token_value" numeric NOT NULL,
  "data" text NULL,
  "chain_id" bigint NOT NULL,
  "block_timestamp" bigint NULL,
  "block_number" bigint NULL,
  "transaction_index" timestamptz NULL,
  PRIMARY KEY ("id", "type")
);
-- create index "idx_transactions_address" to table: "bridge_transactions"
CREATE INDEX "idx_transactions_address" ON "public"."bridge_transactions" ("sender", "receiver");
-- create index "idx_transactions_order" to table: "bridge_transactions"
CREATE INDEX "idx_transactions_order" ON "public"."bridge_transactions" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);
-- create index "idx_transactions_receiver" to table: "bridge_transactions"
CREATE INDEX "idx_transactions_receiver" ON "public"."bridge_transactions" ("receiver");
-- create index "idx_transactions_sender" to table: "bridge_transactions"
CREATE INDEX "idx_transactions_sender" ON "public"."bridge_transactions" ("sender");
-- create "checkpoints" table
CREATE TABLE "public"."checkpoints" (
  "chain_id" bigint NOT NULL,
  "block_number" bigint NOT NULL,
  "block_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("chain_id")
);
-- create "epoch" table
CREATE TABLE "public"."epoch" (
  "id" bigint NOT NULL,
  "start_timestamp" timestamptz NOT NULL,
  "end_timestamp" timestamptz NOT NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NOT NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  "total_operation_rewards" numeric NULL,
  "total_staking_rewards" numeric NULL,
  "total_reward_nodes" bigint NULL,
  "total_request_counts" numeric NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("transaction_hash")
);
-- create index "idx_epoch_id" to table: "epoch"
CREATE INDEX "idx_epoch_id" ON "public"."epoch" ("id" DESC, "block_number" DESC, "transaction_index" DESC);
-- create index "idx_timestamp" to table: "epoch"
CREATE INDEX "idx_timestamp" ON "public"."epoch" ("start_timestamp" DESC, "end_timestamp" DESC);
-- create "epoch_item" table
CREATE TABLE "public"."epoch_item" (
  "epoch_id" bigint NOT NULL,
  "transaction_hash" text NOT NULL,
  "index" bigint NOT NULL,
  "node_address" bytea NOT NULL,
  "operation_rewards" numeric NOT NULL,
  "staking_rewards" numeric NOT NULL,
  "tax_collected" numeric NOT NULL,
  "request_count" numeric NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("transaction_hash", "index")
);
-- create index "idx_epoch_item_epoch_id" to table: "epoch_item"
CREATE INDEX "idx_epoch_item_epoch_id" ON "public"."epoch_item" ("epoch_id");
-- create index "idx_epoch_item_node_address" to table: "epoch_item"
CREATE INDEX "idx_epoch_item_node_address" ON "public"."epoch_item" ("node_address");
-- create "epoch_trigger" table
CREATE TABLE "public"."epoch_trigger" (
  "transaction_hash" text NOT NULL,
  "epoch_id" bigint NOT NULL,
  "data" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("transaction_hash")
);
-- create index "idx_created_at" to table: "epoch_trigger"
CREATE INDEX "idx_created_at" ON "public"."epoch_trigger" ("created_at");
-- create index "idx_epoch_id" to table: "epoch_trigger"
CREATE INDEX "idx_epoch_id" ON "public"."epoch_trigger" ("epoch_id");
-- create "node_count_snapshots" table
CREATE TABLE "public"."node_count_snapshots" (
  "date" date NOT NULL,
  "count" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("date")
);
-- create "node_events" table
CREATE TABLE "public"."node_events" (
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NOT NULL,
  "node_id" bigint NOT NULL,
  "address_from" bytea NOT NULL,
  "address_to" bytea NOT NULL,
  "type" text NOT NULL,
  "log_index" bigint NOT NULL,
  "chain_id" bigint NOT NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  "metadata" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("transaction_hash", "transaction_index", "log_index")
);
-- create index "events_index_address" to table: "node_events"
CREATE INDEX "events_index_address" ON "public"."node_events" ("address_from", "address_to");
-- create index "events_index_address_type" to table: "node_events"
CREATE INDEX "events_index_address_type" ON "public"."node_events" ("address_from", "type");
-- create index "events_index_block_number" to table: "node_events"
CREATE INDEX "events_index_block_number" ON "public"."node_events" ("block_number" DESC, "transaction_index" DESC, "log_index" DESC);
-- create index "events_index_node_id" to table: "node_events"
CREATE INDEX "events_index_node_id" ON "public"."node_events" ("node_id");
-- create "node_info" table
CREATE TABLE "public"."node_info" (
  "address" bytea NOT NULL,
  "id" bigint NOT NULL,
  "endpoint" text NOT NULL,
  "hide_tax_rate" boolean NULL DEFAULT false,
  "is_public_good" boolean NOT NULL,
  "stream" jsonb NULL,
  "config" jsonb NULL,
  "status" text NOT NULL DEFAULT 'offline',
  "last_heartbeat_timestamp" timestamptz NULL,
  "location" jsonb NOT NULL DEFAULT '[]',
  "avatar" jsonb NULL,
  "min_tokens_to_stake" numeric NULL,
  "apy" numeric NULL DEFAULT 0,
  "score" numeric NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("address")
);
-- create index "idx_created_at" to table: "node_info"
CREATE INDEX "idx_created_at" ON "public"."node_info" ("created_at" DESC);
-- create index "idx_endpoint_unique" to table: "node_info"
CREATE UNIQUE INDEX "idx_endpoint_unique" ON "public"."node_info" ("endpoint");
-- create index "idx_id" to table: "node_info"
CREATE UNIQUE INDEX "idx_id" ON "public"."node_info" ("id");
-- create index "idx_is_public" to table: "node_info"
CREATE INDEX "idx_is_public" ON "public"."node_info" ("is_public_good" DESC, "created_at" DESC);
-- create index "idx_last_heartbeat_timestamp" to table: "node_info"
CREATE INDEX "idx_last_heartbeat_timestamp" ON "public"."node_info" ("last_heartbeat_timestamp");
-- create index "idx_score" to table: "node_info"
CREATE INDEX "idx_score" ON "public"."node_info" ("score" DESC);
-- create index "idx_status" to table: "node_info"
CREATE INDEX "idx_status" ON "public"."node_info" ("status");
-- create "node_invalid_response" table
CREATE TABLE "public"."node_invalid_response" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "epoch_id" bigint NOT NULL,
  "type" text NOT NULL,
  "request" text NOT NULL,
  "validator_nodes" bytea[] NULL,
  "validator_response" jsonb NULL,
  "node" bytea NULL,
  "response" jsonb NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create index "idx_epoch_id" to table: "node_invalid_response"
CREATE INDEX "idx_epoch_id" ON "public"."node_invalid_response" ("epoch_id" DESC);
-- create index "idx_node" to table: "node_invalid_response"
CREATE INDEX "idx_node" ON "public"."node_invalid_response" ("node" DESC, "created_at" DESC);
-- create index "idx_request" to table: "node_invalid_response"
CREATE INDEX "idx_request" ON "public"."node_invalid_response" ("request" DESC, "created_at" DESC);
-- create index "idx_type" to table: "node_invalid_response"
CREATE INDEX "idx_type" ON "public"."node_invalid_response" ("type" DESC, "created_at" DESC);
-- create "node_min_tokens_to_stake_snapshots" table
CREATE TABLE "public"."node_min_tokens_to_stake_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NOT NULL,
  "node_address" bytea NOT NULL,
  "epoch_id" bigint NOT NULL,
  "min_tokens_to_stake" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("node_address", "epoch_id")
);
-- create index "min_tokens_to_stake_snapshots_date_idx" to table: "node_min_tokens_to_stake_snapshots"
CREATE INDEX "min_tokens_to_stake_snapshots_date_idx" ON "public"."node_min_tokens_to_stake_snapshots" ("date");
-- create index "min_tokens_to_stake_snapshots_epoch_id_idx" to table: "node_min_tokens_to_stake_snapshots"
CREATE INDEX "min_tokens_to_stake_snapshots_epoch_id_idx" ON "public"."node_min_tokens_to_stake_snapshots" ("epoch_id" DESC, "id" DESC);
-- create index "min_tokens_to_stake_snapshots_id_idx" to table: "node_min_tokens_to_stake_snapshots"
CREATE INDEX "min_tokens_to_stake_snapshots_id_idx" ON "public"."node_min_tokens_to_stake_snapshots" ("id" DESC);
-- create "node_operator_profit_snapshots" table
CREATE TABLE "public"."node_operator_profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NOT NULL,
  "operator" bytea NOT NULL,
  "epoch_id" bigint NOT NULL,
  "operation_pool" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("operator", "epoch_id")
);
-- create index "operator_profit_snapshots_date_idx" to table: "node_operator_profit_snapshots"
CREATE INDEX "operator_profit_snapshots_date_idx" ON "public"."node_operator_profit_snapshots" ("date");
-- create index "operator_profit_snapshots_epoch_id_idx" to table: "node_operator_profit_snapshots"
CREATE INDEX "operator_profit_snapshots_epoch_id_idx" ON "public"."node_operator_profit_snapshots" ("epoch_id" DESC);
-- create index "operator_profit_snapshots_id_idx" to table: "node_operator_profit_snapshots"
CREATE INDEX "operator_profit_snapshots_id_idx" ON "public"."node_operator_profit_snapshots" ("id" DESC);
-- create index "operator_profit_snapshots_operation_pool_idx" to table: "node_operator_profit_snapshots"
CREATE INDEX "operator_profit_snapshots_operation_pool_idx" ON "public"."node_operator_profit_snapshots" ("operation_pool" DESC);
-- create "node_stat" table
CREATE TABLE "public"."node_stat" (
  "address" bytea NOT NULL,
  "endpoint" text NOT NULL,
  "points" numeric NOT NULL,
  "is_public_good" boolean NOT NULL,
  "is_full_node" boolean NOT NULL,
  "is_rss_node" boolean NOT NULL,
  "staking" numeric NOT NULL,
  "epoch" bigint NOT NULL,
  "total_request_count" bigint NOT NULL,
  "epoch_request_count" bigint NOT NULL,
  "epoch_invalid_request_count" bigint NOT NULL,
  "decentralized_network_count" bigint NOT NULL,
  "federated_network_count" bigint NOT NULL,
  "indexer_count" bigint NOT NULL,
  "reset_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("address")
);
-- create index "idx_indexes_created_at" to table: "node_stat"
CREATE INDEX "idx_indexes_created_at" ON "public"."node_stat" ("created_at");
-- create index "idx_indexes_epoch_invalid_request_count" to table: "node_stat"
CREATE INDEX "idx_indexes_epoch_invalid_request_count" ON "public"."node_stat" ("epoch_invalid_request_count");
-- create index "idx_indexes_is_full_node" to table: "node_stat"
CREATE INDEX "idx_indexes_is_full_node" ON "public"."node_stat" ("is_full_node" DESC, "points" DESC);
-- create index "idx_indexes_is_rss_node" to table: "node_stat"
CREATE INDEX "idx_indexes_is_rss_node" ON "public"."node_stat" ("is_rss_node" DESC, "points" DESC);
-- create index "idx_indexes_points" to table: "node_stat"
CREATE INDEX "idx_indexes_points" ON "public"."node_stat" ("points" DESC);
-- create "node_worker" table
CREATE TABLE "public"."node_worker" (
  "address" bytea NOT NULL,
  "network" text NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("address", "network", "name")
);
-- create "stake_chips" table
CREATE TABLE "public"."stake_chips" (
  "id" numeric NOT NULL,
  "owner" text NOT NULL,
  "node" text NOT NULL,
  "value" numeric NULL,
  "metadata" jsonb NULL,
  "block_number" bigint NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_node" to table: "stake_chips"
CREATE INDEX "idx_node" ON "public"."stake_chips" ("node");
-- create index "idx_owner" to table: "stake_chips"
CREATE INDEX "idx_owner" ON "public"."stake_chips" ("owner");
-- create "stake_count_snapshots" table
CREATE TABLE "public"."stake_count_snapshots" (
  "date" date NOT NULL,
  "count" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("date")
);
-- create "stake_events" table
CREATE TABLE "public"."stake_events" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NOT NULL,
  "transaction_status" bigint NOT NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  PRIMARY KEY ("transaction_hash", "block_hash")
);
-- create index "idx_id" to table: "stake_events"
CREATE INDEX "idx_id" ON "public"."stake_events" ("id");
-- create "stake_profit_snapshots" table
CREATE TABLE "public"."stake_profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NOT NULL,
  "owner_address" bytea NOT NULL,
  "epoch_id" bigint NOT NULL,
  "total_chip_amounts" numeric NOT NULL,
  "total_chip_values" numeric NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("owner_address", "epoch_id")
);
-- create index "profit_snapshots_date_idx" to table: "stake_profit_snapshots"
CREATE INDEX "profit_snapshots_date_idx" ON "public"."stake_profit_snapshots" ("date");
-- create index "profit_snapshots_epoch_id_idx" to table: "stake_profit_snapshots"
CREATE INDEX "profit_snapshots_epoch_id_idx" ON "public"."stake_profit_snapshots" ("epoch_id" DESC, "id" DESC);
-- create index "profit_snapshots_id_idx" to table: "stake_profit_snapshots"
CREATE INDEX "profit_snapshots_id_idx" ON "public"."stake_profit_snapshots" ("id" DESC);
-- create index "profit_snapshots_total_chip_amounts_idx" to table: "stake_profit_snapshots"
CREATE INDEX "profit_snapshots_total_chip_amounts_idx" ON "public"."stake_profit_snapshots" ("total_chip_amounts" DESC);
-- create index "profit_snapshots_total_chip_values_idx" to table: "stake_profit_snapshots"
CREATE INDEX "profit_snapshots_total_chip_values_idx" ON "public"."stake_profit_snapshots" ("total_chip_values" DESC);
-- create "stake_transactions" table
CREATE TABLE "public"."stake_transactions" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "user" text NOT NULL,
  "node" text NOT NULL,
  "value" numeric NOT NULL,
  "chips" bigint[] NOT NULL,
  "block_timestamp" timestamptz NOT NULL,
  "block_number" bigint NOT NULL,
  "transaction_index" bigint NOT NULL,
  PRIMARY KEY ("id", "type")
);
-- create index "idx_transactions_address" to table: "stake_transactions"
CREATE INDEX "idx_transactions_address" ON "public"."stake_transactions" ("user", "node");
-- create index "idx_transactions_node" to table: "stake_transactions"
CREATE INDEX "idx_transactions_node" ON "public"."stake_transactions" ("node");
-- create index "idx_transactions_order" to table: "stake_transactions"
CREATE INDEX "idx_transactions_order" ON "public"."stake_transactions" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);
-- create index "idx_transactions_user" to table: "stake_transactions"
CREATE INDEX "idx_transactions_user" ON "public"."stake_transactions" ("user");

-- +goose Down
-- reverse: create index "idx_transactions_user" to table: "stake_transactions"
DROP INDEX "public"."idx_transactions_user";
-- reverse: create index "idx_transactions_order" to table: "stake_transactions"
DROP INDEX "public"."idx_transactions_order";
-- reverse: create index "idx_transactions_node" to table: "stake_transactions"
DROP INDEX "public"."idx_transactions_node";
-- reverse: create index "idx_transactions_address" to table: "stake_transactions"
DROP INDEX "public"."idx_transactions_address";
-- reverse: create "stake_transactions" table
DROP TABLE "public"."stake_transactions";
-- reverse: create index "profit_snapshots_total_chip_values_idx" to table: "stake_profit_snapshots"
DROP INDEX "public"."profit_snapshots_total_chip_values_idx";
-- reverse: create index "profit_snapshots_total_chip_amounts_idx" to table: "stake_profit_snapshots"
DROP INDEX "public"."profit_snapshots_total_chip_amounts_idx";
-- reverse: create index "profit_snapshots_id_idx" to table: "stake_profit_snapshots"
DROP INDEX "public"."profit_snapshots_id_idx";
-- reverse: create index "profit_snapshots_epoch_id_idx" to table: "stake_profit_snapshots"
DROP INDEX "public"."profit_snapshots_epoch_id_idx";
-- reverse: create index "profit_snapshots_date_idx" to table: "stake_profit_snapshots"
DROP INDEX "public"."profit_snapshots_date_idx";
-- reverse: create "stake_profit_snapshots" table
DROP TABLE "public"."stake_profit_snapshots";
-- reverse: create index "idx_id" to table: "stake_events"
DROP INDEX "public"."idx_id";
-- reverse: create "stake_events" table
DROP TABLE "public"."stake_events";
-- reverse: create "stake_count_snapshots" table
DROP TABLE "public"."stake_count_snapshots";
-- reverse: create index "idx_owner" to table: "stake_chips"
DROP INDEX "public"."idx_owner";
-- reverse: create index "idx_node" to table: "stake_chips"
DROP INDEX "public"."idx_node";
-- reverse: create "stake_chips" table
DROP TABLE "public"."stake_chips";
-- reverse: create "node_worker" table
DROP TABLE "public"."node_worker";
-- reverse: create index "idx_indexes_points" to table: "node_stat"
DROP INDEX "public"."idx_indexes_points";
-- reverse: create index "idx_indexes_is_rss_node" to table: "node_stat"
DROP INDEX "public"."idx_indexes_is_rss_node";
-- reverse: create index "idx_indexes_is_full_node" to table: "node_stat"
DROP INDEX "public"."idx_indexes_is_full_node";
-- reverse: create index "idx_indexes_epoch_invalid_request_count" to table: "node_stat"
DROP INDEX "public"."idx_indexes_epoch_invalid_request_count";
-- reverse: create index "idx_indexes_created_at" to table: "node_stat"
DROP INDEX "public"."idx_indexes_created_at";
-- reverse: create "node_stat" table
DROP TABLE "public"."node_stat";
-- reverse: create index "operator_profit_snapshots_operation_pool_idx" to table: "node_operator_profit_snapshots"
DROP INDEX "public"."operator_profit_snapshots_operation_pool_idx";
-- reverse: create index "operator_profit_snapshots_id_idx" to table: "node_operator_profit_snapshots"
DROP INDEX "public"."operator_profit_snapshots_id_idx";
-- reverse: create index "operator_profit_snapshots_epoch_id_idx" to table: "node_operator_profit_snapshots"
DROP INDEX "public"."operator_profit_snapshots_epoch_id_idx";
-- reverse: create index "operator_profit_snapshots_date_idx" to table: "node_operator_profit_snapshots"
DROP INDEX "public"."operator_profit_snapshots_date_idx";
-- reverse: create "node_operator_profit_snapshots" table
DROP TABLE "public"."node_operator_profit_snapshots";
-- reverse: create index "min_tokens_to_stake_snapshots_id_idx" to table: "node_min_tokens_to_stake_snapshots"
DROP INDEX "public"."min_tokens_to_stake_snapshots_id_idx";
-- reverse: create index "min_tokens_to_stake_snapshots_epoch_id_idx" to table: "node_min_tokens_to_stake_snapshots"
DROP INDEX "public"."min_tokens_to_stake_snapshots_epoch_id_idx";
-- reverse: create index "min_tokens_to_stake_snapshots_date_idx" to table: "node_min_tokens_to_stake_snapshots"
DROP INDEX "public"."min_tokens_to_stake_snapshots_date_idx";
-- reverse: create "node_min_tokens_to_stake_snapshots" table
DROP TABLE "public"."node_min_tokens_to_stake_snapshots";
-- reverse: create index "idx_type" to table: "node_invalid_response"
DROP INDEX "public"."idx_type";
-- reverse: create index "idx_request" to table: "node_invalid_response"
DROP INDEX "public"."idx_request";
-- reverse: create index "idx_node" to table: "node_invalid_response"
DROP INDEX "public"."idx_node";
-- reverse: create index "idx_epoch_id" to table: "node_invalid_response"
DROP INDEX "public"."idx_epoch_id";
-- reverse: create "node_invalid_response" table
DROP TABLE "public"."node_invalid_response";
-- reverse: create index "idx_status" to table: "node_info"
DROP INDEX "public"."idx_status";
-- reverse: create index "idx_score" to table: "node_info"
DROP INDEX "public"."idx_score";
-- reverse: create index "idx_last_heartbeat_timestamp" to table: "node_info"
DROP INDEX "public"."idx_last_heartbeat_timestamp";
-- reverse: create index "idx_is_public" to table: "node_info"
DROP INDEX "public"."idx_is_public";
-- reverse: create index "idx_id" to table: "node_info"
DROP INDEX "public"."idx_id";
-- reverse: create index "idx_endpoint_unique" to table: "node_info"
DROP INDEX "public"."idx_endpoint_unique";
-- reverse: create index "idx_created_at" to table: "node_info"
DROP INDEX "public"."idx_created_at";
-- reverse: create "node_info" table
DROP TABLE "public"."node_info";
-- reverse: create index "events_index_node_id" to table: "node_events"
DROP INDEX "public"."events_index_node_id";
-- reverse: create index "events_index_block_number" to table: "node_events"
DROP INDEX "public"."events_index_block_number";
-- reverse: create index "events_index_address_type" to table: "node_events"
DROP INDEX "public"."events_index_address_type";
-- reverse: create index "events_index_address" to table: "node_events"
DROP INDEX "public"."events_index_address";
-- reverse: create "node_events" table
DROP TABLE "public"."node_events";
-- reverse: create "node_count_snapshots" table
DROP TABLE "public"."node_count_snapshots";
-- reverse: create index "idx_epoch_id" to table: "epoch_trigger"
DROP INDEX "public"."idx_epoch_id";
-- reverse: create index "idx_created_at" to table: "epoch_trigger"
DROP INDEX "public"."idx_created_at";
-- reverse: create "epoch_trigger" table
DROP TABLE "public"."epoch_trigger";
-- reverse: create index "idx_epoch_item_node_address" to table: "epoch_item"
DROP INDEX "public"."idx_epoch_item_node_address";
-- reverse: create index "idx_epoch_item_epoch_id" to table: "epoch_item"
DROP INDEX "public"."idx_epoch_item_epoch_id";
-- reverse: create "epoch_item" table
DROP TABLE "public"."epoch_item";
-- reverse: create index "idx_timestamp" to table: "epoch"
DROP INDEX "public"."idx_timestamp";
-- reverse: create index "idx_epoch_id" to table: "epoch"
DROP INDEX "public"."idx_epoch_id";
-- reverse: create "epoch" table
DROP TABLE "public"."epoch";
-- reverse: create "checkpoints" table
DROP TABLE "public"."checkpoints";
-- reverse: create index "idx_transactions_sender" to table: "bridge_transactions"
DROP INDEX "public"."idx_transactions_sender";
-- reverse: create index "idx_transactions_receiver" to table: "bridge_transactions"
DROP INDEX "public"."idx_transactions_receiver";
-- reverse: create index "idx_transactions_order" to table: "bridge_transactions"
DROP INDEX "public"."idx_transactions_order";
-- reverse: create index "idx_transactions_address" to table: "bridge_transactions"
DROP INDEX "public"."idx_transactions_address";
-- reverse: create "bridge_transactions" table
DROP TABLE "public"."bridge_transactions";
-- reverse: create index "idx_id" to table: "bridge_events"
DROP INDEX "public"."idx_id";
-- reverse: create "bridge_events" table
DROP TABLE "public"."bridge_events";
-- reverse: create index "average_tax_rate_submissions_transaction_hash_idx" to table: "average_tax_rate_submissions"
DROP INDEX "public"."average_tax_rate_submissions_transaction_hash_idx";
-- reverse: create index "average_tax_rate_submissions_id_idx" to table: "average_tax_rate_submissions"
DROP INDEX "public"."average_tax_rate_submissions_id_idx";
-- reverse: create "average_tax_rate_submissions" table
DROP TABLE "public"."average_tax_rate_submissions";
