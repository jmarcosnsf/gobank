-- Write your migrate up statements here

CREATE TYPE holder_type AS ENUM ('individual', 'company');
CREATE TYPE account_status AS ENUM ('active', 'closed');
CREATE TYPE transaction_type AS ENUM (
    'deposit',
    'withdrawal',
    'transfer_in',
    'transfer_out'
);

---- create above / drop below ----

DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS account_status;
DROP TYPE IF EXISTS holder_type;


-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
