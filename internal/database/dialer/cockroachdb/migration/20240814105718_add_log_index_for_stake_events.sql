-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."events"
    ADD COLUMN "log_index" integer NOT NULL DEFAULT 0;

ALTER TABLE "stake"."events"
    ADD COLUMN "metadata" jsonb;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."events" DROP COLUMN "metadata";

ALTER TABLE "stake"."events" DROP COLUMN "log_index";
-- +goose StatementEnd
