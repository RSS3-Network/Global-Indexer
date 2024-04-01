-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."epoch_item"
    ADD "request_counts" decimal default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."epoch_item"
    DROP COLUMN "request_counts";
-- +goose StatementEnd
