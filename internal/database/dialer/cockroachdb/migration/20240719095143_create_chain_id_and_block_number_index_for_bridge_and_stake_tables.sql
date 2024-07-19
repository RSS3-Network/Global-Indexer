-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_transactions_chain_id_block_number ON bridge.transactions (chain_id, block_number);
CREATE INDEX idx_events_chain_id_block_number ON bridge.events (chain_id, block_number);

CREATE INDEX idx_transactions_block_number ON stake.transactions (block_number);
CREATE INDEX idx_events_block_number ON stake.events (block_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX bridge.transactions@idx_transactions_chain_id_block_number;
DROP INDEX bridge.events@idx_events_chain_id_block_number;

DROP INDEX stake.transactions@idx_transactions_block_number;
DROP INDEX stake.events@idx_events_block_number;
-- +goose StatementEnd
