-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."epoch"
    ADD "total_request_counts" decimal default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."epoch"
    DROP COLUMN "total_request_counts";
-- +goose StatementEnd
