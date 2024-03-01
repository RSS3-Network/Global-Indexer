-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."epoch_item"
    DROP COLUMN "request_fees";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."epoch_item"
    ADD "request_fees" decimal;
-- +goose StatementEnd
