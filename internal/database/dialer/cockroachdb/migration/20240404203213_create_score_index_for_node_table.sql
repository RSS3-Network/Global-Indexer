-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS "idx_score" ON "public"."node_info" ("score" DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "public"."idx_score";
-- +goose StatementEnd
