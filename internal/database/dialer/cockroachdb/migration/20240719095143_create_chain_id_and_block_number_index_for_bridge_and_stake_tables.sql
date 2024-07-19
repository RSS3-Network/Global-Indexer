-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_transactions_chain_id_block_number ON bridge.transactions (chain_id, block_number) STORING (sender, receiver, token_address_l1, token_address_l2, token_value, data, transaction_index, block_timestamp, finalized);
CREATE INDEX idx_events_chain_id_block_number ON bridge.events (chain_id, block_number) STORING (id, type, transaction_index, transaction_status, block_timestamp, finalized);

CREATE INDEX idx_transactions_block_number ON stake.transactions (block_number) STORING ("user", node, value, chips, transaction_index, block_timestamp, finalized);
CREATE INDEX idx_events_block_number ON stake.events (block_number) STORING (id, type, transaction_index, transaction_status, block_timestamp, finalized);
DROP INDEX stake.chips@idx_block_number;
CREATE INDEX idx_chips_block_number ON stake.chips (block_number) STORING (owner, node, block_timestamp, metadata, value, finalized);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX bridge.transactions@idx_transactions_chain_id_block_number;
DROP INDEX bridge.events@idx_events_chain_id_block_number;

DROP INDEX stake.transactions@idx_transactions_block_number;
DROP INDEX stake.events@idx_events_block_number;
CREATE INDEX idx_block_number ON stake.chips (block_number);
DROP INDEX stake.chips@idx_chips_block_number;
-- +goose StatementEnd
