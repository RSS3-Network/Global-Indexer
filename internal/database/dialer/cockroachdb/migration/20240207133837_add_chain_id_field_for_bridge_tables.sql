-- +goose Up
-- +goose StatementBegin
-- Add the chain id field for the transactions table.
ALTER TABLE "bridge"."transactions"
    ADD "chain_id" bigint;

UPDATE "bridge"."transactions"
SET "chain_id" = 11155111
WHERE "type" = 'deposit';

UPDATE "bridge"."transactions"
SET "chain_id" = 2331
WHERE "type" = 'withdraw';

-- Add the chain id field for the events table.
ALTER TABLE "bridge"."events"
    ADD "chain_id" bigint;

UPDATE "bridge"."events"
SET "chain_id" = "transactions"."chain_id"
FROM "bridge"."transactions"
WHERE "events"."id" = "transactions"."id";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "bridge"."transactions" DROP COLUMN "chain_id";
ALTER TABLE "bridge"."events" DROP COLUMN "chain_id";
-- +goose StatementEnd
