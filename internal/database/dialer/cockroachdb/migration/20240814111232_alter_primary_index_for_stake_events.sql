-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."events" DROP CONSTRAINT IF EXISTS "pk_stake_events";

ALTER TABLE "stake"."events" ADD CONSTRAINT "pk_stake" PRIMARY KEY ("transaction_hash", "log_index", "id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."events" ADD CONSTRAINT "pk_stake_events" PRIMARY KEY ("transaction_hash", "block_hash", "id");

ALTER TABLE "stake"."events" DROP CONSTRAINT "pk_stake";
-- +goose StatementEnd
