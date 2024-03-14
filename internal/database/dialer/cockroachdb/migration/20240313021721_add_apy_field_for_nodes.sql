-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    ADD "apy" decimal DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    DROP COLUMN "apy";
-- +goose StatementEnd
