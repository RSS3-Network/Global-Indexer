-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE node_info
    ADD COLUMN name text NOT NULL DEFAULT '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE node_info
    DROP COLUMN name;