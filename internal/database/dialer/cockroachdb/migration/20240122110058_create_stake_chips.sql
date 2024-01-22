-- +goose Up
-- +goose StatementBegin
CREATE TABLE "stake"."chips"
(
    "id"    decimal NOT NULL UNIQUE,
    "owner" text    NOT NULL,
    "node"  text    NOT NULL,

    CONSTRAINT "pk_chips" PRIMARY KEY ("id")
);

CREATE INDEX "idx_owner" ON "stake"."chips" ("owner");
CREATE INDEX "idx_node" ON "stake"."chips" ("node");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "stake"."chips";
-- +goose StatementEnd
