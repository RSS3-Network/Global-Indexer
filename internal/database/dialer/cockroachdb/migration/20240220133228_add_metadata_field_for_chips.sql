-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    ADD "metadata" jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    DROP COLUMN "metadata";
-- +goose StatementEnd
