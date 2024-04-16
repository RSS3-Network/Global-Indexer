-- +goose Up
-- create "average_tax_rate_submissions" table
CREATE TABLE "public"."average_tax_rate_submissions" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "epoch_id" bigint NULL,
  "average_tax_rate" text NULL,
  "transaction_hash" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "bridge_events" table
CREATE TABLE "public"."bridge_events" (
  "id" text NULL,
  "type" text NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NULL,
  "transaction_status" bigint NULL,
  "chain_id" bigint NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NULL,
  "block_timestamp" timestamptz NULL,
  PRIMARY KEY ("transaction_hash", "block_hash")
);
-- create "bridge_transactions" table
CREATE TABLE "public"."bridge_transactions" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "sender" text NULL,
  "receiver" text NULL,
  "token_address_l1" text NULL,
  "token_address_l2" text NULL,
  "token_value" text NULL,
  "data" text NULL,
  "chain_id" bigint NULL,
  "block_timestamp" timestamptz NULL,
  "block_number" bigint NULL,
  "transaction_index" bigint NULL,
  PRIMARY KEY ("id", "type")
);
-- create "checkpoints" table
CREATE TABLE "public"."checkpoints" (
  "chain_id" bigint NULL,
  "block_number" bigint NULL,
  "block_hash" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "rowid" bigint NOT NULL DEFAULT unique_rowid(),
  PRIMARY KEY ("rowid")
);
-- create "epoch" table
CREATE TABLE "public"."epoch" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "start_timestamp" timestamptz NULL,
  "end_timestamp" timestamptz NULL,
  "transaction_hash" text NULL,
  "transaction_index" bigint NULL,
  "block_hash" text NULL,
  "block_number" bigint NULL,
  "block_timestamp" timestamptz NULL,
  "total_operation_rewards" text NULL,
  "total_staking_rewards" text NULL,
  "total_reward_items" bigint NULL,
  "total_request_counts" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "epoch_item" table
CREATE TABLE "public"."epoch_item" (
  "epoch_id" bigint NULL,
  "index" bigint NOT NULL,
  "transaction_hash" text NOT NULL,
  "node_address" text NULL,
  "operation_rewards" text NULL,
  "staking_rewards" text NULL,
  "tax_amounts" text NULL,
  "request_counts" text NULL,
  PRIMARY KEY ("index", "transaction_hash")
);
-- create "epoch_trigger" table
CREATE TABLE "public"."epoch_trigger" (
  "transaction_hash" text NULL,
  "epoch_id" bigint NULL,
  "data" bytea NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "rowid" bigint NOT NULL DEFAULT unique_rowid(),
  PRIMARY KEY ("rowid")
);
-- create "node_count_snapshots" table
CREATE TABLE "public"."node_count_snapshots" (
  "date" timestamptz NULL,
  "count" bigint NULL,
  "rowid" bigint NOT NULL DEFAULT unique_rowid(),
  PRIMARY KEY ("rowid")
);
-- create "node_events" table
CREATE TABLE "public"."node_events" (
  "transaction_hash" text NULL,
  "transaction_index" bigint NULL,
  "node_id" bigint NULL,
  "address_from" bytea NULL,
  "address_to" bytea NULL,
  "type" text NULL,
  "log_index" bigint NULL,
  "chain_id" bigint NULL,
  "block_hash" text NULL,
  "block_number" bigint NULL,
  "block_timestamp" timestamptz NULL,
  "metadata" bytea NULL,
  "rowid" bigint NOT NULL DEFAULT unique_rowid(),
  PRIMARY KEY ("rowid")
);
-- create "node_indexer" table
CREATE TABLE "public"."node_indexer" (
  "address" bytea NOT NULL,
  "network" text NOT NULL,
  "worker" text NOT NULL,
  PRIMARY KEY ("address", "network", "worker")
);
-- create "node_info" table
CREATE TABLE "public"."node_info" (
  "address" bytea NOT NULL,
  "id" bigint NULL,
  "endpoint" text NULL,
  "hide_tax_rate" boolean NULL,
  "is_public_good" boolean NULL,
  "stream" bytea NULL,
  "config" jsonb NULL,
  "status" text NULL,
  "last_heartbeat_timestamp" timestamptz NULL,
  "local" jsonb NULL,
  "avatar" jsonb NULL,
  "min_tokens_to_stake" text NULL,
  "apy" text NULL,
  "score" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("address")
);
-- create "node_min_tokens_to_stake_snapshots" table
CREATE TABLE "public"."node_min_tokens_to_stake_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "node_address" bytea NULL,
  "min_tokens_to_stake" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "node_operator_profit_snapshots" table
CREATE TABLE "public"."node_operator_profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "operator" bytea NULL,
  "operation_pool" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "node_stat" table
CREATE TABLE "public"."node_stat" (
  "address" bytea NOT NULL,
  "endpoint" text NULL,
  "points" numeric NULL,
  "is_public_good" boolean NULL,
  "is_full_node" boolean NULL,
  "is_rss_node" boolean NULL,
  "staking" numeric NULL,
  "epoch" bigint NULL,
  "total_request_count" bigint NULL,
  "epoch_request_count" bigint NULL,
  "epoch_invalid_request_count" bigint NULL,
  "decentralized_network_count" bigint NULL,
  "federated_network_count" bigint NULL,
  "indexer_count" bigint NULL,
  "reset_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("address")
);
-- create "stake_chips" table
CREATE TABLE "public"."stake_chips" (
  "id" text NOT NULL,
  "owner" text NULL,
  "node" text NULL,
  "value" text NULL,
  "metadata" bytea NULL,
  "block_number" text NULL,
  "block_timestamp" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "stake_count_snapshots" table
CREATE TABLE "public"."stake_count_snapshots" (
  "date" timestamptz NULL,
  "count" bigint NULL,
  "rowid" bigint NOT NULL DEFAULT unique_rowid(),
  PRIMARY KEY ("rowid")
);
-- create "stake_events" table
CREATE TABLE "public"."stake_events" (
  "id" text NULL,
  "type" text NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NULL,
  "transaction_status" bigint NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NULL,
  "block_timestamp" timestamptz NULL,
  PRIMARY KEY ("transaction_hash", "block_hash")
);
-- create "stake_profit_snapshots" table
CREATE TABLE "public"."stake_profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "owner_address" bytea NULL,
  "total_chip_amounts" text NULL,
  "total_chip_values" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "stake_transactions" table
CREATE TABLE "public"."stake_transactions" (
  "id" text NOT NULL,
  "type" text NOT NULL,
  "user" text NULL,
  "node" text NULL,
  "value" text NULL,
  "chips" bigint[] NULL,
  "block_timestamp" timestamptz NULL,
  "block_number" bigint NULL,
  "transaction_index" bigint NULL,
  PRIMARY KEY ("id", "type")
);

-- +goose Down
-- reverse: create "stake_transactions" table
DROP TABLE "public"."stake_transactions";
-- reverse: create "stake_profit_snapshots" table
DROP TABLE "public"."stake_profit_snapshots";
-- reverse: create "stake_events" table
DROP TABLE "public"."stake_events";
-- reverse: create "stake_count_snapshots" table
DROP TABLE "public"."stake_count_snapshots";
-- reverse: create "stake_chips" table
DROP TABLE "public"."stake_chips";
-- reverse: create "node_stat" table
DROP TABLE "public"."node_stat";
-- reverse: create "node_operator_profit_snapshots" table
DROP TABLE "public"."node_operator_profit_snapshots";
-- reverse: create "node_min_tokens_to_stake_snapshots" table
DROP TABLE "public"."node_min_tokens_to_stake_snapshots";
-- reverse: create "node_info" table
DROP TABLE "public"."node_info";
-- reverse: create "node_indexer" table
DROP TABLE "public"."node_indexer";
-- reverse: create "node_events" table
DROP TABLE "public"."node_events";
-- reverse: create "node_count_snapshots" table
DROP TABLE "public"."node_count_snapshots";
-- reverse: create "epoch_trigger" table
DROP TABLE "public"."epoch_trigger";
-- reverse: create "epoch_item" table
DROP TABLE "public"."epoch_item";
-- reverse: create "epoch" table
DROP TABLE "public"."epoch";
-- reverse: create "checkpoints" table
DROP TABLE "public"."checkpoints";
-- reverse: create "bridge_transactions" table
DROP TABLE "public"."bridge_transactions";
-- reverse: create "bridge_events" table
DROP TABLE "public"."bridge_events";
-- reverse: create "average_tax_rate_submissions" table
DROP TABLE "public"."average_tax_rate_submissions";
