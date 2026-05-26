-- Write your migrate up statements here

CREATE TABLE transaction (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id              UUID NOT NULL REFERENCES account(id),
    type                    transaction_type NOT NULL,
    amount                  NUMERIC(15,2) NOT NULL CHECK (amount > 0),
    counterparty_account_id UUID REFERENCES account(id),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transaction_account_id ON transaction(account_id, created_at DESC);

---- create above / drop below ----

DROP TABLE IF EXISTS transaction;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
