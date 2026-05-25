package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmarcosnsf/gobank/internal/store/pgstore"
	"github.com/shopspring/decimal"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountNotOwned = errors.New("account does not belong to the authenticated holder")
)

type AccountService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewAccountService(pool *pgxpool.Pool) AccountService {
	return AccountService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func accountBelongsTo(account pgstore.GetAccountByIDRow, holderID uuid.UUID, holderType HolderType) bool {
	if holderType == HolderTypeIndividual {
		return account.IndividualHolderID.Valid && account.IndividualHolderID.Bytes == holderID
	}
	return account.CompanyHolderID.Valid && account.CompanyHolderID.Bytes == holderID
}

func (as *AccountService) CreateAccount(ctx context.Context, holderID uuid.UUID, holderType HolderType) (uuid.UUID, error) {
	args := pgstore.CreateAccountParams{
		HolderType: pgstore.HolderType(holderType),
	}

	if holderType == HolderTypeIndividual {
		args.IndividualHolderID = pgtype.UUID{Bytes: holderID, Valid: true}
	} else {
		args.CompanyHolderID = pgtype.UUID{Bytes: holderID, Valid: true}
	}

	id, err := as.queries.CreateAccount(ctx, args)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (as *AccountService) GetBalance(ctx context.Context, accountID, holderID uuid.UUID, holderType HolderType) (decimal.Decimal, error){
	account, err := as.queries.GetAccountByID(ctx, accountID)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return decimal.Decimal{}, ErrAccountNotFound
		}
		return decimal.Decimal{}, err
	}

	if !accountBelongsTo(account, holderID, holderType){
		return decimal.Decimal{}, ErrAccountNotOwned
	}

	return account.Balance, nil
}