-- +goose Up
-- +goose StatementBegin
ALTER TABLE "stake"."profit_snapshots"
    RENAME COLUMN "total_chip_amounts" TO "total_chip_amount";
ALTER TABLE "stake"."profit_snapshots"
    RENAME COLUMN "total_chip_values" TO "total_chip_value";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "stake"."profit_snapshots"
    RENAME COLUMN "total_chip_amount" TO "total_chip_amounts";
ALTER TABLE "stake"."profit_snapshots"
    RENAME COLUMN "total_chip_value" TO "total_chip_values";
-- +goose StatementEnd
