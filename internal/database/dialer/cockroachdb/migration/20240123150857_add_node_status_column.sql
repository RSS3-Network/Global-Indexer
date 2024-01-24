-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_info ADD COLUMN status TEXT;
ALTER TABLE node_info ADD COLUMN last_heartbeat_timestamp timestamptz;
CREATE INDEX IF NOT EXISTS "idx_status" ON "node_info" ("status");
CREATE INDEX IF NOT EXISTS "idx_last_heartbeat_timestamp" ON "node_info" ("last_heartbeat_timestamp");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_status" ON "node_info";
DROP INDEX IF EXISTS "idx_last_heartbeat_timestamp" ON "node_info";
ALTER TABLE node_info DROP COLUMN status;
ALTER TABLE node_info DROP COLUMN last_heartbeat_timestamp;
-- +goose StatementEnd
