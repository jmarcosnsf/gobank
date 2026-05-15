-- Write your migrate up statements here

CREATE TABLE account (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    individual_holder_id UUID REFERENCES individual_holder(id),
    company_holder_id    UUID REFERENCES company_holder(id),
    holder_type          holder_type NOT NULL,
    balance              NUMERIC(15,2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    status               account_status NOT NULL DEFAULT 'active',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at            TIMESTAMPTZ,
    CHECK (
        (holder_type = 'individual' AND individual_holder_id IS NOT NULL AND company_holder_id IS NULL) OR
        (holder_type = 'company'    AND company_holder_id    IS NOT NULL AND individual_holder_id IS NULL)
    )
);

CREATE INDEX idx_account_individual_holder ON account(individual_holder_id) WHERE individual_holder_id IS NOT NULL;
CREATE INDEX idx_account_company_holder    ON account(company_holder_id)    WHERE company_holder_id    IS NOT NULL;

---- create above / drop below ----

DROP TABLE IF EXISTS account;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
