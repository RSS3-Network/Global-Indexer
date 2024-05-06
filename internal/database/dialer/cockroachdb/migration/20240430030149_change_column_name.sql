-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_info"
    RENAME COLUMN "local" TO "location";
ALTER TABLE "epoch"
    RENAME COLUMN "total_reward_items" TO "total_reward_nodes";
ALTER TABLE "epoch_item"
    RENAME COLUMN "tax_amounts" TO "tax_collected";
ALTER TABLE "epoch_item"
    RENAME COLUMN "request_counts" TO "request_count";
ALTER TABLE "epoch_item"
    RENAME TO "node_reward_record";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_reward_record"
    RENAME TO "epoch_item";
ALTER TABLE "epoch_item"
    RENAME COLUMN "request_count" TO "request_counts";
ALTER TABLE "epoch_item"
    RENAME COLUMN "tax_collected" TO "tax_amounts";
ALTER TABLE "epoch"
    RENAME COLUMN "total_reward_nodes" TO "total_reward_items";
ALTER TABLE "node_info"
    RENAME COLUMN "location" TO "local";
-- +goose StatementEnd
