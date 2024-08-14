-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS "idx_events_order" ON "stake"."events" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "stake"."idx_events_order";
-- +goose StatementEnd
