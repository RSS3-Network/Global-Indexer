-- +goose Up
-- +goose StatementBegin
ALTER TABLE "bridge"."transactions"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE "bridge"."events"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE "stake"."transactions"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE "stake"."events"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE "stake"."chips"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE "node"."events"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "bridge"."transactions"
    DROP COLUMN "finalized";
ALTER TABLE "bridge"."events"
    DROP COLUMN "finalized";

ALTER TABLE "stake"."transactions"
    DROP COLUMN "finalized";
ALTER TABLE "stake"."events"
    DROP COLUMN "finalized";
ALTER TABLE "stake"."chips"
    DROP COLUMN "finalized";

ALTER TABLE "node"."events"
    DROP COLUMN "finalized";
-- +goose StatementEnd
