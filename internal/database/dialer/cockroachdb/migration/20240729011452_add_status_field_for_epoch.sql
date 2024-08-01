-- +goose Up
-- +goose StatementBegin
ALTER TABLE "epoch"
    ADD COLUMN "finalized" BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "epoch"
    DROP COLUMN "finalized";
-- +goose StatementEnd
