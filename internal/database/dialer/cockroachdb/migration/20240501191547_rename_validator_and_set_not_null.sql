-- +goose Up
-- +goose StatementBegin
ALTER TABLE node_invalid_response
    ALTER COLUMN validator_nodes DROP NOT NULL;

ALTER TABLE node_invalid_response
    ALTER COLUMN validator_response DROP NOT NULL;

ALTER TABLE "node_invalid_response"
    RENAME COLUMN "validator_nodes" TO "verifier_nodes";

ALTER TABLE "node_invalid_response"
    RENAME COLUMN "validator_response" TO "verifier_response";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_invalid_response"
    RENAME COLUMN "verifier_nodes" TO "validator_nodes";

ALTER TABLE "node_invalid_response"
    RENAME COLUMN "verifier_response" TO "validator_response";

ALTER TABLE node_invalid_response
    ALTER COLUMN validator_nodes SET NOT NULL;

ALTER TABLE node_invalid_response
    ALTER COLUMN validator_response SET NOT NULL;
-- +goose StatementEnd
