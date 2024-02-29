-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    ADD "min_tokens_to_stake" decimal;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."chips"
    DROP COLUMN "min_tokens_to_stake";
-- +goose StatementEnd
