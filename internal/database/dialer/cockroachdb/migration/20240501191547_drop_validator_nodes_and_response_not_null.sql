-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_invalid_response
    ALTER COLUMN validator_nodes DROP NOT NULL;

ALTER TABLE node_invalid_response
    ALTER COLUMN validator_response DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE node_invalid_response
    ALTER COLUMN validator_nodes SET NOT NULL;

ALTER TABLE node_invalid_response
    ALTER COLUMN validator_response SET NOT NULL;
-- +goose StatementEnd
