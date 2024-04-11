-- Create "average_tax_rate_submissions" table
CREATE TABLE "public"."average_tax_rate_submissions" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "epoch_id" bigint NULL,
  "average_tax_rate" text NULL,
  "transaction_hash" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_average_tax_rate_submissions_deleted_at" to table: "average_tax_rate_submissions"
CREATE INDEX "idx_average_tax_rate_submissions_deleted_at" ON "public"."average_tax_rate_submissions" ("deleted_at");
-- Create "checkpoints" table
CREATE TABLE "public"."checkpoints" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "chain_id" bigint NULL,
  "block_number" bigint NULL,
  "block_hash" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_checkpoints_deleted_at" to table: "checkpoints"
CREATE INDEX "idx_checkpoints_deleted_at" ON "public"."checkpoints" ("deleted_at");
-- Create "chips" table
CREATE TABLE "public"."chips" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "owner" text NULL,
  "node" text NULL,
  "value" text NULL,
  "metadata" bytea NULL,
  "block_number" text NULL,
  "block_timestamp" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_chips_deleted_at" to table: "chips"
CREATE INDEX "idx_chips_deleted_at" ON "public"."chips" ("deleted_at");
-- Create "count_snapshots" table
CREATE TABLE "public"."count_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "date" timestamptz NULL,
  "count" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_count_snapshots_deleted_at" to table: "count_snapshots"
CREATE INDEX "idx_count_snapshots_deleted_at" ON "public"."count_snapshots" ("deleted_at");
-- Create "epoch" table
CREATE TABLE "public"."epoch" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
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
  PRIMARY KEY ("id")
);
-- Create index "idx_epoch_deleted_at" to table: "epoch"
CREATE INDEX "idx_epoch_deleted_at" ON "public"."epoch" ("deleted_at");
-- Create "epoch_item" table
CREATE TABLE "public"."epoch_item" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "epoch_id" bigint NULL,
  "index" bigint NOT NULL,
  "transaction_hash" text NOT NULL,
  "node_address" text NULL,
  "operation_rewards" text NULL,
  "staking_rewards" text NULL,
  "tax_amounts" text NULL,
  "request_counts" text NULL,
  PRIMARY KEY ("id", "index", "transaction_hash")
);
-- Create index "idx_epoch_item_deleted_at" to table: "epoch_item"
CREATE INDEX "idx_epoch_item_deleted_at" ON "public"."epoch_item" ("deleted_at");
-- Create "epoch_trigger" table
CREATE TABLE "public"."epoch_trigger" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "transaction_hash" text NULL,
  "epoch_id" bigint NULL,
  "data" bytea NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_epoch_trigger_deleted_at" to table: "epoch_trigger"
CREATE INDEX "idx_epoch_trigger_deleted_at" ON "public"."epoch_trigger" ("deleted_at");
-- Create "events" table
CREATE TABLE "public"."events" (
  "id" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "type" text NULL,
  "transaction_hash" text NOT NULL,
  "transaction_index" bigint NULL,
  "transaction_status" bigint NULL,
  "block_hash" text NOT NULL,
  "block_number" bigint NULL,
  "block_timestamp" timestamptz NULL,
  PRIMARY KEY ("transaction_hash", "block_hash")
);
-- Create index "idx_events_deleted_at" to table: "events"
CREATE INDEX "idx_events_deleted_at" ON "public"."events" ("deleted_at");
-- Create "min_tokens_to_stake_snapshots" table
CREATE TABLE "public"."min_tokens_to_stake_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "node_address" bytea NULL,
  "min_tokens_to_stake" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_min_tokens_to_stake_snapshots_deleted_at" to table: "min_tokens_to_stake_snapshots"
CREATE INDEX "idx_min_tokens_to_stake_snapshots_deleted_at" ON "public"."min_tokens_to_stake_snapshots" ("deleted_at");
-- Create "node_indexer" table
CREATE TABLE "public"."node_indexer" (
  "address" bytea NOT NULL,
  "network" text NOT NULL,
  "worker" text NOT NULL,
  PRIMARY KEY ("address", "network", "worker")
);
-- Create "node_info" table
CREATE TABLE "public"."node_info" (
  "id" bigint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "address" bytea NOT NULL,
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
  PRIMARY KEY ("address")
);
-- Create index "idx_node_info_deleted_at" to table: "node_info"
CREATE INDEX "idx_node_info_deleted_at" ON "public"."node_info" ("deleted_at");
-- Create "node_stat" table
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
-- Create "operator_profit_snapshots" table
CREATE TABLE "public"."operator_profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "operator" bytea NULL,
  "operation_pool" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_operator_profit_snapshots_deleted_at" to table: "operator_profit_snapshots"
CREATE INDEX "idx_operator_profit_snapshots_deleted_at" ON "public"."operator_profit_snapshots" ("deleted_at");
-- Create "profit_snapshots" table
CREATE TABLE "public"."profit_snapshots" (
  "id" bigint NOT NULL DEFAULT unique_rowid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "date" timestamptz NULL,
  "epoch_id" bigint NULL,
  "owner_address" bytea NULL,
  "total_chip_amounts" text NULL,
  "total_chip_values" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_profit_snapshots_deleted_at" to table: "profit_snapshots"
CREATE INDEX "idx_profit_snapshots_deleted_at" ON "public"."profit_snapshots" ("deleted_at");
-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
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
-- Create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "public"."transactions" ("deleted_at");
