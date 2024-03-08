-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    ADD "hide_tax_rate" bool DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    DROP COLUMN "hide_tax_rate";
-- +goose StatementEnd
