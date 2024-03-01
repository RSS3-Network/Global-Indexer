-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    ADD "min_tokens_to_stake" decimal;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    DROP COLUMN "min_tokens_to_stake";
-- +goose StatementEnd
