package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmarcosnsf/gobank/internal/store/pgstore"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)
var(
	ErrDuplicateEmailOrCpf = errors.New("cpf or email already exists")
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