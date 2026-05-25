package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmarcosnsf/gobank/internal/store/pgstore"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)
var(
	ErrDuplicateEmailOrCpf = errors.New("cpf or email already exists")
	ErrDuplicateEmailOrCnpj = errors.New("cnpj or email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type HolderService struct {
	pool *pgxpool.Pool
	queries *pgstore.Queries
}

func NewHolderService(pool *pgxpool.Pool) HolderService {
	return HolderService{
		pool: pool,
		queries: pgstore.New(pool),
	}
} 

func(hs *HolderService) CreateIndividualHolder(
	ctx context.Context,
	fullName, cpf, email, phone, category string,
	dateOfBirth time.Time,
	monthlyIncome decimal.Decimal,
	password string,
) (uuid.UUID, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateIndividualHolderParams{
		FullName: fullName,
		Cpf: cpf,
		DateOfBirth: pgtype.Date{Time: dateOfBirth, Valid: true},
		Email: email,
		Phone: phone,
		Category: pgtype.Text{String: category, Valid: true},
		MonthlyIncome: decimal.NullDecimal{Decimal: monthlyIncome, Valid: true},
		PasswordHash: hash,
	}

	id, err := hs.queries.CreateIndividualHolder(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505"{
			return uuid.UUID{}, ErrDuplicateEmailOrCpf
		}
		return uuid.UUID{}, err
	}

	return id, nil
}

type HolderType string

const (
    HolderTypeIndividual HolderType = "individual"
    HolderTypeCompany    HolderType = "company"
)

func(hs *HolderService) Authenticate(ctx context.Context, email, password string) (uuid.UUID, HolderType, error){
	individual, err := hs.queries.GetIndividualHolderByEmail(ctx, email)
	if err == nil {
		if err := bcrypt.CompareHashAndPassword(individual.PasswordHash, []byte(password)); err != nil {
			return uuid.UUID{}, "", ErrInvalidCredentials
		}
		return individual.ID, HolderTypeIndividual, nil
	}
	if !errors.Is(err, pgx.ErrNoRows){
		return uuid.UUID{}, "", err
	}

	company, err := hs.queries.GetCompanyHolderByEmail(ctx, email)
	if err == nil{
		if err := bcrypt.CompareHashAndPassword(company.PasswordHash, []byte(password)); err != nil {
			return uuid.UUID{}, "", ErrInvalidCredentials
		}
		return company.ID, HolderTypeCompany, nil
	}
	if errors.Is(err, pgx.ErrNoRows){
		return uuid.UUID{}, "", ErrInvalidCredentials
	}

	return uuid.UUID{}, "", err
}

func (hs *HolderService) CreateCompanyHolder(
	ctx context.Context,
	tradeName, cnpj, corporateEmail, phone, category string,
	foundedAt time.Time,
	annualRevenue decimal.Decimal,
	password string,
) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.UUID{}, err
	}

	args := pgstore.CreateCompanyHolderParams{
		TradeName:      tradeName,
		Cnpj:           cnpj,
		FoundedAt:      pgtype.Date{Time: foundedAt, Valid: true},
		CorporateEmail: corporateEmail,
		Phone:          phone,
		Category:       pgtype.Text{String: category, Valid: true},
		AnnualRevenue:  decimal.NullDecimal{Decimal: annualRevenue, Valid: true},
		PasswordHash:   hash,
	}

	id, err := hs.queries.CreateCompanyHolder(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, ErrDuplicateEmailOrCnpj
		}
		return uuid.UUID{}, err
	}

	return id, nil
}