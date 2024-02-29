-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    ADD "value" decimal;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    DROP COLUMN "value";
-- +goose StatementEnd
