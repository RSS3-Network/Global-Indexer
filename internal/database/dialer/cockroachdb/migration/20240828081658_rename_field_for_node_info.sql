-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_info RENAME COLUMN "type" TO version;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE node_info RENAME COLUMN version TO "type";
-- +goose StatementEnd
