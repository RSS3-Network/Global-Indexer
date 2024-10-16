-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    DROP COLUMN "score";

DROP INDEX IF EXISTS "public"."idx_score";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."node_info"
    ADD "score" decimal default 0;

CREATE INDEX IF NOT EXISTS "idx_score" ON "public"."node_info" ("score" DESC);
-- +goose StatementEnd