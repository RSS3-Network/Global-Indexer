-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE node_stat
    ADD COLUMN is_rsshub_node BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_node_stat_is_rsshub_node_points
    ON node_stat (is_rsshub_node, points DESC);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS idx_node_stat_is_rsshub_node_points;

ALTER TABLE node_stat
DROP COLUMN is_rsshub_node;