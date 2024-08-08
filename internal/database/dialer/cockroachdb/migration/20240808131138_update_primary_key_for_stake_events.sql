-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."events" DROP CONSTRAINT "pk_events";

ALTER TABLE "stake"."events" ADD CONSTRAINT "pk_stake_events" PRIMARY KEY ("transaction_hash", "block_hash", "id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."events" ADD CONSTRAINT "pk_events" PRIMARY KEY ("transaction_hash", "block_hash");

ALTER TABLE "stake"."events" DROP CONSTRAINT "pk_stake_events";
-- +goose StatementEnd