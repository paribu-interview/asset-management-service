CREATE TABLE IF NOT EXISTS assets
(
    "id"         serial PRIMARY KEY,
    "created_at" timestamp      NOT NULL DEFAULT now(),
    "updated_at" timestamp               DEFAULT NULL,
    "wallet_id"  integer        NOT NULL,
    "name"       VARCHAR(255)   NOT NULL,
    "amount"     NUMERIC(18, 2) NOT NULL DEFAULT 0.0,
    unique (wallet_id, name)
);