-- +migrate Up

create table transactions
(
    id                 bigserial primary key not null,
    balance_wei        numeric(78, 0)       not null,
    sender             varchar(42)          not null,
    recipient          varchar(42)          not null,
    transaction_hash   varchar(66)          not null unique,
    transaction_index  varchar(66)          ,
    timestamp          timestamp without time zone default now()
);

create index transactions_sender_index on transactions (sender);
create index transactions_recipient_index on transactions (recipient);
create index transactions_transaction_hash_index on transactions (transaction_hash);

-- +migrate Down

drop index transactions_transaction_hash_index;
drop index transactions_recipient_index;
drop index transactions_sender_index;
drop table transactions;