-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    ADD "score" decimal default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
DROP COLUMN "score";
-- +goose StatementEnd
