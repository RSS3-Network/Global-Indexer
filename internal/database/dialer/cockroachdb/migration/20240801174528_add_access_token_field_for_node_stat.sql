-- +goose Up
-- +goose StatementBegin
ALTER TABLE "node_stat"
    ADD COLUMN "access_token" TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "node_stat"
DROP COLUMN "access_token";
-- +goose StatementEnd