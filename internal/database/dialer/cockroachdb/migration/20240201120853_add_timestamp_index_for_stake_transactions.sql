-- +goose Up
-- +goose StatementBegin
CREATE TABLE "stake"."transactions_upgraded"
(
    "id"                text        NOT NULL,
    "type"              text        NOT NULL,
    "user"              text        NOT NULL,
    "node"              text        NOT NULL,
    "value"             decimal     NOT NULL,
    "chips"             bigint[]    NOT NULL,
    "block_number"      bigint      NOT NULL,
    "transaction_index" integer     NOT NULL,
    "block_timestamp"   timestamptz NOT NULL,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

INSERT INTO "stake"."transactions_upgraded"
SELECT "transactions".*, "events"."block_number", "events"."transaction_index", "events"."block_timestamp"
FROM "stake"."transactions"
         LEFT JOIN "stake"."events"
                   ON "transactions"."id" = "events"."id" AND (
                       ("transactions"."type" = 'deposit' AND "events"."type" = 'deposited')
                           OR
                       ("transactions"."type" = 'withdraw' AND "events"."type" = 'requested')
                           OR
                       ("transactions"."type" = 'stake' AND "events"."type" = 'staked')
                           OR
                       ("transactions"."type" = 'unstake' AND "events"."type" = 'requested')
                       );

CREATE INDEX "idx_transactions_user" ON "stake"."transactions_upgraded" ("user");
CREATE INDEX "idx_transactions_node" ON "stake"."transactions_upgraded" ("node");
CREATE INDEX "idx_transactions_address" ON "stake"."transactions_upgraded" ("user", "node");
CREATE INDEX "idx_transactions_order" ON "stake"."transactions_upgraded" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);

DROP TABLE "stake"."transactions";

ALTER TABLE "stake"."transactions_upgraded"
    RENAME TO "transactions";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE "stake"."transactions_downgraded"
(
    "id"    text     NOT NULL,
    "type"  text     NOT NULL,
    "user"  text     NOT NULL,
    "node"  text     NOT NULL,
    "value" decimal  NOT NULL,
    "chips" bigint[] NOT NULL,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

INSERT INTO "stake"."transactions_downgraded"
SELECT "id",
       "type",
       "user",
       "node",
       "value",
       "chips"
FROM "stake"."transactions";

DROP TABLE "stake"."transactions";

ALTER TABLE "stake"."transactions_downgraded"
    RENAME TO "transactions";
-- +goose StatementEnd
