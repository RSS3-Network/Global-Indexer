-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS "gateway";

CREATE TABLE gateway.account
(
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    deleted_at   timestamp with time zone,
    address      bytea not null
        primary key,
    ru_limit     bigint,
    is_paused    boolean,
    billing_rate numeric
);

CREATE INDEX idx_gateway_account_deleted_at
    ON gateway.account (deleted_at);

CREATE TABLE gateway.key
(
    id                bigint default unique_rowid() not null
        primary key,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    key               text
        constraint idx_gateway_key_key
            unique,
    ru_used_total     bigint,
    ru_used_current   bigint,
    api_calls_total   bigint,
    api_calls_current bigint,
    name              text,
    account_address   bytea
        constraint fk_gateway_key_account
            references gateway.account
);

CREATE INDEX idx_gateway_key_deleted_at
    ON gateway.key (deleted_at);

CREATE INDEX idx_gateway_key_account_address
    ON gateway.key (account_address);

CREATE TABLE gateway.consumption_log
(
    id               bigint default unique_rowid() not null
        primary key,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    consumption_date timestamp with time zone,
    ru_used          bigint,
    api_calls        bigint,
    key_id           bigint
        constraint fk_gateway_consumption_log_key
            references gateway.key
);

CREATE INDEX idx_gateway_consumption_log_key_id
    ON gateway.consumption_log (key_id);

CREATE INDEX idx_gateway_consumption_log_consumption_date
    ON gateway.consumption_log (consumption_date);

CREATE TABLE gateway.pending_withdraw_request
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    amount          numeric,
    account_address bytea not null
        primary key
        constraint fk_gateway_pending_withdraw_request_account
            references gateway.account
);

CREATE TABLE gateway.br_deposited
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    tx_hash         bytea not null
        primary key,
    "index" bigint,
    block_timestamp timestamp with time zone,
    "user"          bytea,
    amount          text
);

CREATE INDEX idx_gateway_br_deposited_block_timestamp
    ON gateway.br_deposited (block_timestamp);

CREATE TABLE gateway.br_withdrawn
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    tx_hash         bytea not null
        primary key,
    "index" bigint,
    block_timestamp timestamp with time zone,
    "user"          bytea,
    amount          text,
    fee             text
);

CREATE INDEX idx_gateway_br_withdrawn_block_timestamp
    ON gateway.br_withdrawn (block_timestamp);

CREATE TABLE gateway.br_collected
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    tx_hash         bytea not null
        primary key,
    "index" bigint,
    block_timestamp timestamp with time zone,
    "user"          bytea,
    amount          text
);

CREATE INDEX idx_gateway_br_collected_block_timestamp
    ON gateway.br_collected (block_timestamp);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gateway.consumption_log CASCADE;

DROP TABLE IF EXISTS gateway.key CASCADE;

DROP TABLE IF EXISTS gateway.pending_withdraw_request CASCADE;

DROP TABLE IF EXISTS gateway.account CASCADE;

DROP TABLE IF EXISTS gateway.br_deposited CASCADE;

DROP TABLE IF EXISTS gateway.br_withdrawn CASCADE;

DROP TABLE IF EXISTS gateway.br_collected CASCADE;

drop schema if exists gateway CASCADE;
-- +goose StatementEnd
