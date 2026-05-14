-- Write your migrate up statements here

CREATE TABLE individual_holder (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name      VARCHAR(255) NOT NULL,
    cpf            VARCHAR(14)  UNIQUE NOT NULL,
    date_of_birth  DATE         NOT NULL,
    email          VARCHAR(255) UNIQUE NOT NULL,
    phone          VARCHAR(20)  NOT NULL,
    category       VARCHAR(50),
    monthly_income NUMERIC(15,2),
    password_hash  BYTEA        NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE company_holder (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_name       VARCHAR(255) NOT NULL,
    cnpj             VARCHAR(18)  UNIQUE NOT NULL,
    founded_at       DATE         NOT NULL,
    corporate_email  VARCHAR(255) UNIQUE NOT NULL,
    phone            VARCHAR(20)  NOT NULL,
    category         VARCHAR(50),
    annual_revenue   NUMERIC(15,2),
    password_hash    BYTEA        NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

---- create above / drop below ----

DROP TABLE IF EXISTS company_holder;
DROP TABLE IF EXISTS individual_holder;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
