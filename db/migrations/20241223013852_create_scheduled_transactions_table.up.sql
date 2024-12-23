CREATE TABLE IF NOT EXISTS scheduled_transactions
(
    "id"                    serial PRIMARY KEY,
    "created_at"            timestamp      NOT NULL DEFAULT now(),
    "updated_at"            timestamp               DEFAULT NULL,
    "source_wallet_id"      integer        NOT NULL,
    "destination_wallet_id" integer        NOT NULL,
    "asset_name"            VARCHAR(255)   NOT NULL,
    "amount"                NUMERIC(18, 2) NOT NULL,
    "scheduled_at"          timestamp      NOT NULL,
    "status"                VARCHAR(255)   NOT NULL DEFAULT 'pending'
);

CREATE INDEX idx_scheduled_time_status ON scheduled_transactions (scheduled_at, status);