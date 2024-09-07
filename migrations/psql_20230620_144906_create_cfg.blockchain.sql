-- +migrate Up
CREATE TABLE IF NOT EXISTS cfg.blockchain
(
    id               uuid        NOT NULL PRIMARY KEY,
    fee_blion        float4      NOT NULL,
    fee_miner        float4      NOT NULL,
    fee_validator    float4      NOT NULL,
    fee_node         float4      NOT NULL,
    ttl_block        INTEGER     NOT NULL,
    max_transactions INTEGER     NOT NULL,
    max_miners       int4        NOT NULL,
    max_validators   int4        NOT NULL,
    tickets_price    int4        NOT NULL,
    lottery_ttl      int4        NOT NULL,
    wallet_main      varchar(50) NOT NULL,
    deleted_at       TIMESTAMP,
    created_at       TIMESTAMP   NOT NULL DEFAULT now(),
    updated_at       TIMESTAMP   NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS cfg.blockchain;
