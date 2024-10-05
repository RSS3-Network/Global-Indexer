-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_info RENAME COLUMN "version" TO type;
ALTER TABLE "node_info" ADD COLUMN "version" TEXT;
CREATE INDEX IF NOT EXISTS "idx_node_info_type" ON "node_info" ("type");
CREATE INDEX IF NOT EXISTS "idx_node_info_version" ON "node_info" ("version");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_node_info_type";
DROP INDEX IF EXISTS "idx_node_info_version";
ALTER TABLE "node_info" DROP COLUMN IF EXISTS "version";
ALTER TABLE node_info RENAME COLUMN type TO "version";
CREATE INDEX IF NOT EXISTS "idx_node_info_version" ON "node_info" ("version");
-- +goose StatementEnd