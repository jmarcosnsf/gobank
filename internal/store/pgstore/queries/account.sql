-- name: CreateAccount :one
INSERT INTO account (individual_holder_id, company_holder_id, holder_type)
VALUES ($1, $2, $3)
RETURNING id;