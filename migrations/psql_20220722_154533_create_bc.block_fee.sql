
-- +migrate Up
CREATE TABLE IF NOT EXISTS bc.block_fee(
    id uuid NOT NULL PRIMARY KEY,
    block_id BIGINT  NOT NULL,
    fee float8  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS bc.block_fee;
