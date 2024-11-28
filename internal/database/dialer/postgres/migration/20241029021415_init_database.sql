-- +goose Up
-- +goose StatementBegin

-- bridge schema
create schema if not exists "bridge";

create table if not exists "bridge"."events"
(
    id                 text                     not null,
    type               text                     not null,
    transaction_hash   text                     not null,
    transaction_index  bigint                   not null,
    transaction_status bigint                   not null,
    chain_id           bigint                   not null,
    block_hash         text                     not null,
    block_number       bigint                   not null,
    block_timestamp    timestamp with time zone not null,
    finalized          boolean default false    not null,
    constraint pk_bridge_events primary key (transaction_hash, block_hash)
);

create index if not exists "idx_bridge_events_id" on "bridge"."events" (id);

create index if not exists "idx_bridge_events_chain_id_block_number" on "bridge"."events" (chain_id, block_number);

create table if not exists "bridge"."transactions"
(
    id                text                  not null,
    type              text                  not null,
    sender            text                  not null,
    receiver          text                  not null,
    token_address_l1  text,
    token_address_l2  text,
    token_value       numeric               not null,
    data              text,
    chain_id          bigint                not null,
    block_number      bigint,
    transaction_index bigint,
    block_timestamp   timestamp with time zone,
    finalized         boolean default false not null,
    constraint pk_bridge_transactions primary key (id, type)
);

create index if not exists "idx_bridge_transactions_sender" on "bridge"."transactions" (sender);

create index if not exists "idx_bridge_transactions_sender_receiver" on "bridge"."transactions" (sender, receiver);

create index if not exists "idx_bridge_transactions_receiver" on "bridge"."transactions" (receiver);

create index if not exists "idx_bridge_transactions_chain_id_block_number" on "bridge"."transactions" (chain_id, block_number);

create index if not exists "idx_bridge_transactions_order" on "bridge"."transactions" (block_timestamp desc, block_number desc, transaction_index desc);

-- epoch schema
create schema if not exists "epoch";

create table if not exists "epoch"."apy_snapshots"
(
    epoch_id   bigint                                 not null,
    date       timestamp with time zone               not null,
    apy        numeric                                not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    constraint pk_epoch_apy_snapshots primary key (epoch_id)
);

create index if not exists "idx_epoch_apy_snapshots_date" on "epoch"."apy_snapshots" (date);

-- node schema
create schema if not exists "node";

create table if not exists "node"."count_snapshots"
(
    date  date             not null,
    count bigint default 0 not null,
    constraint pk_node_count_snapshots primary key (date)
);

create table if not exists "node"."events"
(
    transaction_hash  text                                   not null,
    transaction_index bigint                                 not null,
    node_id           bigint                                 not null,
    address_from      bytea                                  not null,
    address_to        bytea                                  not null,
    type              text                                   not null,
    log_index         bigint                                 not null,
    chain_id          bigint                                 not null,
    block_hash        text                                   not null,
    block_number      bigint                                 not null,
    block_timestamp   timestamp with time zone               not null,
    metadata          jsonb                                  not null,
    created_at        timestamp with time zone default now() not null,
    updated_at        timestamp with time zone default now() not null,
    finalized         boolean                  default false not null,
    constraint pk_node_events primary key (transaction_hash, transaction_index, log_index)
);

create index if not exists "idx_node_events_node_id" on "node"."events" (node_id);

create index if not exists "idx_node_events_address_from_address_to" on "node"."events" (address_from, address_to);

create index if not exists "idx_node_events_address_from_type" on "node"."events" (address_from, type);

create index if not exists "idx_node_events_block_number_transaction_index_log_index" on "node"."events" (block_number desc, transaction_index desc, log_index desc);

create index if not exists "idx_node_events_block_number" on "node"."events" (block_number);

create sequence if not exists "node"."operator_profit_snapshots_id_seq" minvalue 0;

create table if not exists "node"."operator_profit_snapshots"
(
    id             bigint                   default nextval('"node"."operator_profit_snapshots_id_seq"'::REGCLASS) not null,
    date           timestamp with time zone                                                                        not null,
    epoch_id       bigint                                                                                          not null,
    operator       bytea                                                                                           not null,
    operation_pool numeric                                                                                         not null,
    created_at     timestamp with time zone default now()                                                          not null,
    updated_at     timestamp with time zone default now()                                                          not null,
    constraint pk_node_operator_profit_snapshots primary key (operator, epoch_id)
);

create index if not exists "idx_node_operator_profit_snapshots_operation_pool" on "node"."operator_profit_snapshots" (operation_pool desc);

create index if not exists "idx_node_operator_profit_snapshots_epoch_id" on "node"."operator_profit_snapshots" (epoch_id desc);

create index if not exists "idx_node_operator_profit_snapshots_id" on "node"."operator_profit_snapshots" (id desc);

create index if not exists "idx_node_operator_profit_snapshots_date" on "node"."operator_profit_snapshots" (date);

create sequence if not exists "node"."apy_snapshots_id_seq" minvalue 0;

create table if not exists "node"."apy_snapshots"
(
    id           bigint                   default nextval('"node"."apy_snapshots_id_seq"'::REGCLASS) not null,
    date         timestamp with time zone                                                            not null,
    epoch_id     bigint                                                                              not null,
    node_address bytea                                                                               not null,
    apy          numeric                                                                             not null,
    created_at   timestamp with time zone default now()                                              not null,
    updated_at   timestamp with time zone default now()                                              not null,
    constraint pk_node_apy_snapshots primary key (node_address, epoch_id)
);

create index if not exists "idx_node_apy_snapshots_date" on "node"."apy_snapshots" (date);

create index if not exists "idx_node_apy_snapshots_epoch_id_id" on "node"."apy_snapshots" (epoch_id desc, id desc);

-- stake schema
create schema if not exists "stake";

create table if not exists "stake"."transactions"
(
    id                text                     not null,
    type              text                     not null,
    "user"            text                     not null,
    node              text                     not null,
    value             numeric                  not null,
    chips             bigint[]                 not null,
    block_number      bigint                   not null,
    transaction_index bigint                   not null,
    block_timestamp   timestamp with time zone not null,
    finalized         boolean default false    not null,
    constraint pk_stake_transactions primary key (id, type)
);

create index if not exists "idx_stake_transactions_block_number" on "stake"."transactions" (block_number);

create index if not exists "idx_stake_transactions_order" on "stake"."transactions" (block_timestamp desc, block_number desc, transaction_index desc);

create index if not exists "idx_stake_transactions_node" on "stake"."transactions" (node);

create index if not exists "idx_stake_transactions_user_node" on "stake"."transactions" ("user", node);

create index if not exists "idx_stake_transactions_user" on "stake"."transactions" ("user");

create table if not exists "stake"."events"
(
    id                 text                     not null,
    type               text                     not null,
    transaction_hash   text                     not null,
    transaction_index  bigint                   not null,
    transaction_status bigint                   not null,
    block_hash         text                     not null,
    block_number       bigint                   not null,
    block_timestamp    timestamp with time zone not null,
    finalized          boolean default false    not null,
    log_index          bigint  default 0        not null,
    metadata           jsonb,
    constraint pk_stake_events primary key (transaction_hash, log_index, id)
);

create index if not exists "idx_stake_events_id" on "stake"."events" (id);

create index if not exists "idx_stake_events_order" on "stake"."events" (block_timestamp desc, block_number desc, transaction_index desc);

create index if not exists "idx_stake_events_block_number" on "stake"."events" (block_number);

create table if not exists "stake"."chips"
(
    id              numeric                  not null,
    owner           text                     not null,
    node            text                     not null,
    block_number    bigint                   not null,
    block_timestamp timestamp with time zone not null,
    metadata        jsonb,
    value           numeric,
    finalized       boolean default false    not null,
    constraint pk_stake_chips primary key (id)
);

create index if not exists "idx_stake_chips_owner_node_value_finalized" on "stake"."chips" (owner, node, value, finalized);

create index if not exists "idx_stake_chips_block_number" on "stake"."chips" (block_number);

create index if not exists "idx_stake_chips_node" on "stake"."chips" (node);

create index if not exists "idx_stake_chips_owner" on "stake"."chips" (owner);

create table if not exists "stake"."count_snapshots"
(
    date  date             not null,
    count bigint default 0 not null,
    constraint pk_stake_count_snapshots primary key (date)
);

create sequence if not exists "stake"."profit_snapshots_id_seq" minvalue 0;

create table if not exists "stake"."profit_snapshots"
(
    id                 bigint                   default nextval('"stake"."profit_snapshots_id_seq"'::REGCLASS) not null,
    date               timestamp with time zone                                                                not null,
    epoch_id           bigint                                                                                  not null,
    owner_address      bytea                                                                                   not null,
    total_chip_amounts numeric                                                                                 not null,
    total_chip_values  numeric                                                                                 not null,
    created_at         timestamp with time zone default now()                                                  not null,
    updated_at         timestamp with time zone default now()                                                  not null,
    constraint pk_stake_profit_snapshots primary key (owner_address, epoch_id)
);

create index if not exists "idx_stake_profit_snapshots_date" on "stake"."profit_snapshots" (date);

create index if not exists "idx_stake_profit_snapshots_total_chip_amounts" on "stake"."profit_snapshots" (total_chip_amounts desc);

create index if not exists "idx_stake_profit_snapshots_epoch_id_id" on "stake"."profit_snapshots" (epoch_id desc, id desc);

create index if not exists "idx_stake_profit_snapshots_id" on "stake"."profit_snapshots" (id);

create index if not exists "idx_stake_profit_snapshots_total_chip_values" on "stake"."profit_snapshots" (total_chip_values desc);

create view "stake"."stakings"(staker, node, count, value) as
SELECT owner AS staker, node, count(*) AS count, sum(value) AS value
FROM "stake"."chips"
WHERE finalized IS true
GROUP BY owner, node;

-- public
create sequence if not exists "average_tax_rate_submissions_id_seq" minvalue 0;

create sequence if not exists "node_invalid_response_id_seq" minvalue 0;

create table if not exists "node_info"
(
    id                       bigint                                           not null,
    address                  bytea                                            not null,
    endpoint                 text                                             not null,
    is_public_good           boolean                                          not null,
    stream                   jsonb,
    config                   jsonb,
    status                   text                     default 'offline'::text not null,
    location                 jsonb                    default '[]'::JSONB     not null,
    last_heartbeat_timestamp timestamp with time zone,
    created_at               timestamp with time zone default now()           not null,
    updated_at               timestamp with time zone default now()           not null,
    avatar                   jsonb,
    hide_tax_rate            boolean                  default false,
    apy                      numeric                  default 0,
    type                     text,
    access_token             text,
    version                  text,
    constraint pk_node_info primary key (address),
    constraint idx_id unique (id),
    constraint idx_endpoint_unique unique (endpoint)
);

create index if not exists "idx_node_info_is_public" on "node_info" (is_public_good asc, created_at desc);

create index if not exists "idx_node_info_last_heartbeat_timestamp" on "node_info" (last_heartbeat_timestamp);

create index if not exists "idx_node_info_status" on "node_info" (status);

create index if not exists "idx_node_info_address_created_at" on "node_info" (address asc, created_at desc);

create index if not exists "idx_node_info_version" on "node_info" (version);

create index if not exists "idx_node_info_type" on "node_info" (type);

create table if not exists "checkpoints"
(
    chain_id     bigint                                 not null,
    block_number bigint                                 not null,
    block_hash   text                                   not null,
    created_at   timestamp with time zone default now() not null,
    updated_at   timestamp with time zone default now() not null,
    constraint pk_checkpoints primary key (chain_id)
);

create table if not exists "node_stat"
(
    address                     bytea                                  not null,
    endpoint                    text                                   not null,
    points                      numeric                                not null,
    is_public_good              boolean                                not null,
    is_full_node                boolean                                not null,
    is_rss_node                 boolean                                not null,
    staking                     numeric                                not null,
    epoch                       bigint                                 not null,
    total_request_count         bigint                                 not null,
    epoch_request_count         bigint                                 not null,
    epoch_invalid_request_count bigint                                 not null,
    decentralized_network_count bigint                                 not null,
    federated_network_count     bigint                                 not null,
    indexer_count               bigint                                 not null,
    reset_at                    timestamp with time zone               not null,
    created_at                  timestamp with time zone default now() not null,
    updated_at                  timestamp with time zone default now() not null,
    access_token                text,
    constraint pk_node_stat primary key (address)
);

create index if not exists "idx_node_stat_created_at" on "node_stat" (created_at);

create index if not exists "idx_node_stat_epoch_invalid_request_count" on "node_stat" (epoch_invalid_request_count);

create index if not exists "idx_node_stat_points" on "node_stat" (points desc);

create index if not exists "idx_node_stat_is_full_node_points" on "node_stat" (is_full_node, points desc);

create index if not exists "idx_node_stat_is_rss_node_points" on "node_stat" (is_rss_node, points desc);

create table if not exists "average_tax_rate_submissions"
(
    id               bigint                   default nextval('"average_tax_rate_submissions_id_seq"'::REGCLASS) not null,
    epoch_id         bigint                                                                                      not null,
    transaction_hash text                                                                                        not null,
    average_tax_rate numeric                                                                                     not null,
    created_at       timestamp with time zone default now()                                                      not null,
    updated_at       timestamp with time zone default now()                                                      not null,
    constraint pk_average_tax_rate_submissions primary key (epoch_id)
);

create index if not exists "idx_average_tax_rate_submissions_transaction_hash" on "average_tax_rate_submissions" (transaction_hash);

create index if not exists "idx_average_tax_rate_submissions_id" on "average_tax_rate_submissions" (id desc);

create index if not exists "idx_average_tax_rate_submissions_epoch_id" on "average_tax_rate_submissions" (epoch_id desc);

create table if not exists "epoch"
(
    id                      bigint                                 not null,
    start_timestamp         timestamp with time zone               not null,
    end_timestamp           timestamp with time zone               not null,
    block_hash              text                                   not null,
    block_number            bigint                                 not null,
    block_timestamp         timestamp with time zone               not null,
    transaction_hash        text                                   not null,
    transaction_index       bigint                                 not null,
    total_operation_rewards numeric,
    total_staking_rewards   numeric,
    total_rewarded_nodes    bigint,
    created_at              timestamp with time zone default now() not null,
    updated_at              timestamp with time zone default now() not null,
    total_request_counts    numeric                  default 0,
    finalized               boolean                  default false not null,
    constraint pk_epoch primary key (transaction_hash)
);

create index if not exists "idx_epoch_start_timestamp_end_timestamp" on "epoch" (start_timestamp desc, end_timestamp desc);

create index if not exists "idx_epoch_id_block_number_transaction_index" on "epoch" (id desc, block_number desc, transaction_index desc);

create table if not exists "epoch_trigger"
(
    transaction_hash text                                   not null,
    epoch_id         bigint                                 not null,
    data             jsonb                                  not null,
    created_at       timestamp with time zone default now() not null,
    updated_at       timestamp with time zone default now() not null,
    constraint pk_epoch_trigger primary key (transaction_hash)
);

create index if not exists "idx_epoch_trigger_epoch_id" on "epoch_trigger" (epoch_id);

create index if not exists "idx_epoch_trigger_created_at" on "epoch_trigger" (created_at);

create table if not exists "node_invalid_response"
(
    id                bigint                   default nextval('"node_invalid_response_id_seq"'::REGCLASS) not null,
    epoch_id          bigint                                                                               not null,
    type              text                                                                                 not null,
    request           text                                                                                 not null,
    verifier_nodes    bytea[],
    verifier_response jsonb,
    node              bytea                                                                                not null,
    response          jsonb                                                                                not null,
    created_at        timestamp with time zone default now()                                               not null,
    updated_at        timestamp with time zone default now()                                               not null,
    constraint pk_node_invalid_response primary key (id)
);

create index if not exists "idx_node_invalid_response_node_created_at" on "node_invalid_response" (node asc, created_at desc);

create index if not exists "idx_node_invalid_response_request_created_at" on "node_invalid_response" (request asc, created_at desc);

create index if not exists "idx_node_invalid_response_type_created_at" on "node_invalid_response" (type asc, created_at desc);

create index if not exists "idx_node_invalid_response_epoch_id" on "node_invalid_response" (epoch_id desc);

create table if not exists "node_reward_record"
(
    epoch_id          bigint                                 not null,
    index             bigint                                 not null,
    node_address      bytea                                  not null,
    transaction_hash  text                                   not null,
    operation_rewards numeric                                not null,
    staking_rewards   numeric                                not null,
    tax_collected     numeric                                not null,
    created_at        timestamp with time zone default now() not null,
    updated_at        timestamp with time zone default now() not null,
    request_count     numeric                  default 0,
    constraint pk_epoch_item primary key (transaction_hash, index)
);

create index if not exists "idx_node_reward_record_node_address" on "node_reward_record" (node_address);

create index if not exists "idx_node_reward_record_epoch_id" on "node_reward_record" (epoch_id);

create table if not exists "node_worker"
(
    address   bytea                 not null,
    network   text                  not null,
    name      text                  not null,
    epoch_id  bigint  default 0     not null,
    is_active boolean default false not null,
    constraint pk_node_worker primary key (epoch_id, address, network, name)
);

create index if not exists "idx_node_worker_is_active" on "node_worker" (is_active);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "bridge"."events";
drop table if exists "bridge"."transactions";
drop schema if exists "bridge";

drop table if exists "epoch"."apy_snapshots";
drop schema if exists "epoch";

drop table if exists "node"."events";
drop sequence if exists "node"."operator_profit_snapshots_id_seq";
drop table if exists "node"."operator_profit_snapshots";
drop sequence if exists "node"."apy_snapshots_id_seq";
drop table if exists "node"."apy_snapshots";
drop table if exists "node"."count_snapshots";
drop schema if exists "node";

drop table if exists "stake"."transactions";
drop table if exists "stake"."events";
drop table if exists "stake"."chips";
drop sequence if exists "stake"."profit_snapshots_id_seq";
drop table if exists "stake"."profit_snapshots";
drop table if exists "stake"."count_snapshots";
drop schema if exists "stake";

drop table if exists "checkpoints";
drop table if exists "node_info";
drop table if exists "node_stat";
drop table if exists "node_worker";
drop table if exists "epoch";
drop table if exists "node_reward_record";
drop table if exists "epoch_trigger";
drop sequence if exists "average_tax_rate_submissions_id_seq";
drop table if exists "average_tax_rate_submissions";
drop sequence if exists "node_invalid_response_id_seq";
drop table if exists "node_invalid_response";
-- +goose StatementEnd
