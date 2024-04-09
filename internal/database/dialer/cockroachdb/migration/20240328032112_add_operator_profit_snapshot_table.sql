-- +goose Up
-- +goose StatementBegin
CREATE INDEX profit_snapshots_id_idx ON "stake"."profit_snapshots" ("id");

CREATE INDEX "min_tokens_to_stake_snapshots_id_idx" ON "node"."min_tokens_to_stake_snapshots" ("id" DESC);

CREATE TABLE "node"."operator_profit_snapshots"
(
    "id"             bigint GENERATED BY DEFAULT AS IDENTITY (INCREMENT 1 MINVALUE 0 START 0),
    "date"           timestamptz NOT NULL,
    "epoch_id"       bigint      NOT NULL,
    "operator"       bytea       NOT NULL,
    "operation_pool" decimal     NOT NULL,
    "created_at"     timestamptz NOT NULL DEFAULT now(),
    "updated_at"     timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT "pkey" PRIMARY KEY ("operator", "epoch_id")
);

CREATE INDEX "operator_profit_snapshots_date_idx" ON "node"."operator_profit_snapshots" ("date");
CREATE INDEX "operator_profit_snapshots_epoch_id_idx" ON "node"."operator_profit_snapshots" ("epoch_id" DESC);
CREATE INDEX "operator_profit_snapshots_id_idx" ON "node"."operator_profit_snapshots" ("id" DESC);
CREATE INDEX "operator_profit_snapshots_operation_pool_idx" ON "node"."operator_profit_snapshots" ("operation_pool" DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "stake"."profit_snapshots_id_idx";
DROP INDEX "node"."min_tokens_to_stake_snapshots_id_idx";
DROP TABLE "node"."operator_profit_snapshots";
-- +goose StatementEnd
