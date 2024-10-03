-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_node_info_version";

ALTER TABLE node_info RENAME COLUMN "version" TO type;

ALTER TABLE "node_info" ADD COLUMN "version" TEXT;

CREATE INDEX "idx_node_info_type" ON "node_info" ("type");

CREATE INDEX "idx_node_info_version" ON "node_info" ("version");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_node_info_version";
DROP INDEX IF EXISTS "idx_node_info_type";

ALTER TABLE "node_info" DROP COLUMN "version";
ALTER TABLE node_info RENAME COLUMN type TO "version";

CREATE INDEX "idx_node_info_version" ON "node_info" ("version");
-- +goose StatementEnd