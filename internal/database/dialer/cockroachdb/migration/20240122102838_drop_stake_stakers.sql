-- +goose Up
-- +goose StatementBegin
DROP TABLE "stake"."stakers";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE "stake"."stakers"
(
    "user"  text   NOT NULL,
    "node"  text   NOT NULL,
    "value" decimal NOT NULL,

    CONSTRAINT "pk_stakers" PRIMARY KEY ("user", "node")
);

CREATE INDEX "idx_value" ON "stake"."stakers" ("value");
-- +goose StatementEnd
