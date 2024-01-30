-- +goose Up
-- +goose StatementBegin
CREATE TABLE "bridge"."transactions_upgraded"
(
    "id"                text    NOT NULL,
    "type"              text    NOT NULL,
    "sender"            text    NOT NULL,
    "receiver"          text    NOT NULL,
    "token_address_l1"  text,
    "token_address_l2"  text,
    "token_value"       decimal NOT NULL,
    "data"              text,
    "block_number"      bigint,
    "transaction_index" integer,
    "block_timestamp"   timestamptz,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

INSERT INTO "bridge"."transactions_upgraded"
SELECT "transactions".*, "events"."block_number", "events"."transaction_index", "events"."block_timestamp"
FROM "bridge"."transactions"
         LEFT JOIN "bridge"."events"
                   ON "transactions"."id" = "events"."id" AND "events"."type" = 'initialized';

CREATE INDEX "idx_transactions_sender" ON "bridge"."transactions_upgraded" ("sender");
CREATE INDEX "idx_transactions_receiver" ON "bridge"."transactions_upgraded" ("receiver");
CREATE INDEX "idx_transactions_address" ON "bridge"."transactions_upgraded" ("sender", "receiver");
CREATE INDEX "idx_transactions_order" ON "bridge"."transactions_upgraded" ("block_timestamp" DESC, "block_number" DESC, "transaction_index" DESC);

DROP TABLE "bridge"."transactions";

ALTER TABLE "bridge"."transactions_upgraded"
    RENAME TO "transactions";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE "bridge"."transactions_downgraded"
(
    "id"               text    NOT NULL,
    "type"             text    NOT NULL,
    "sender"           text    NOT NULL,
    "receiver"         text    NOT NULL,
    "token_address_l1" text,
    "token_address_l2" text,
    "token_value"      decimal NOT NULL,
    "data"             text,

    CONSTRAINT "pk_transactions" PRIMARY KEY ("id", "type")
);

INSERT INTO "bridge"."transactions_downgraded"
SELECT "id",
       "type",
       "sender",
       "receiver",
       "token_address_l1",
       "token_address_l2",
       "token_value",
       "data"
FROM "bridge"."transactions";

DROP TABLE "bridge"."transactions";

ALTER TABLE "bridge"."transactions_downgraded"
    RENAME TO "transactions";
-- +goose StatementEnd
