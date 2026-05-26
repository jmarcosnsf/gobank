-- name: CreateAccount :one
INSERT INTO account (individual_holder_id, company_holder_id, holder_type)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetAccountByID :one
SELECT id, individual_holder_id, company_holder_id, holder_type, balance, status
FROM account
WHERE id = $1;

-- name: GetAccountByIDForUpdate :one
SELECT id, balance, status 
FROM account
WHERE id = $1
FOR UPDATE;

-- name: UpdateAccountBalance :exec
UPDATE account SET balance = $2 WHERE id = $1;

-- name: CreateTransaction :one
INSERT INTO transaction (account_id, type, amount, counterparty_account_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CloseAccount :exec
UPDATE account
SET status = 'closed', closed_at = NOW()
WHERE id = $1;