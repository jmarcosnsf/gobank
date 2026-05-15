-- name: CreateIndividualHolder :one
INSERT INTO individual_holder (full_name, cpf, date_of_birth, email, phone, category, monthly_income, password_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: CreateCompanyHolder :one
INSERT INTO company_holder (trade_name, cnpj, founded_at, corporate_email, phone, category, annual_revenue, password_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;