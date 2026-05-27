package services

import (
	"bytes"
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
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountNotOwned   = errors.New("account does not belong to the authenticated holder")
	ErrAccountNotActive  = errors.New("account is not active")
	ErrAccountHasBalance = errors.New("account still has balance")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameAccountTransfer = errors.New("cannot transfer to the same account")
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

func (as *AccountService) GetBalance(ctx context.Context, accountID, holderID uuid.UUID, holderType HolderType) (decimal.Decimal, error) {
	account, err := as.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return decimal.Decimal{}, ErrAccountNotFound
		}
		return decimal.Decimal{}, err
	}

	if !accountBelongsTo(account, holderID, holderType) {
		return decimal.Decimal{}, ErrAccountNotOwned
	}

	return account.Balance, nil
}

func (as *AccountService) Deposit(
	ctx context.Context,
	accountID, holderID uuid.UUID,
	holderType HolderType,
	amount decimal.Decimal,
) error {
	account, err := as.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}
	if !accountBelongsTo(account, holderID, holderType) {
		return ErrAccountNotOwned
	}

	tx, err := as.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := as.queries.WithTx(tx)

	locked, err := qtx.GetAccountByIDForUpdate(ctx, accountID)
	if err != nil {
		return err
	}
	if locked.Status != pgstore.AccountStatusActive {
		return ErrAccountNotActive
	}

	newBalance := locked.Balance.Add(amount)

	err = qtx.UpdateAccountBalance(ctx, pgstore.UpdateAccountBalanceParams{
		ID:      accountID,
		Balance: newBalance,
	})
	if err != nil {
		return err
	}

	_, err = qtx.CreateTransaction(ctx, pgstore.CreateTransactionParams{
		AccountID: accountID,
		Type:      pgstore.TransactionTypeDeposit,
		Amount:    amount,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (as *AccountService) CloseAccount(
	ctx context.Context,
	accountID, holderID uuid.UUID,
	holderType HolderType,
) error {
	account, err := as.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}
	if !accountBelongsTo(account, holderID, holderType) {
		return ErrAccountNotOwned
	}

	tx, err := as.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := as.queries.WithTx(tx)

	locked, err := qtx.GetAccountByIDForUpdate(ctx, accountID)
	if err != nil {
		return err
	}

	if locked.Status != pgstore.AccountStatusActive {
		return ErrAccountNotActive
	}
	if !locked.Balance.IsZero() {
		return ErrAccountHasBalance
	}

	if err := qtx.CloseAccount(ctx, accountID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (as *AccountService) Withdraw(
	ctx context.Context,
	accountID, holderID uuid.UUID,
	holderType HolderType,
	amount decimal.Decimal,
) error {
	account, err := as.queries.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}
	if !accountBelongsTo(account, holderID, holderType) {
		return ErrAccountNotOwned
	}

	tx, err := as.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := as.queries.WithTx(tx)

	locked, err := qtx.GetAccountByIDForUpdate(ctx, accountID)
	if err != nil {
		return err
	}
	if locked.Status != pgstore.AccountStatusActive {
		return ErrAccountNotActive
	}
	if locked.Balance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	newBalance := locked.Balance.Sub(amount)

	err = qtx.UpdateAccountBalance(ctx, pgstore.UpdateAccountBalanceParams{
		ID:      accountID,
		Balance: newBalance,
	})
	if err != nil {
		return err
	}

	_, err = qtx.CreateTransaction(ctx, pgstore.CreateTransactionParams{
		AccountID: accountID,
		Type:      pgstore.TransactionTypeWithdrawal,
		Amount:    amount,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (as *AccountService) Transfer(
	ctx context.Context,
	fromAccountID, toAccountID, holderID uuid.UUID,
	holderType HolderType,
	amount decimal.Decimal,
) error {
	if fromAccountID == toAccountID {
		return ErrSameAccountTransfer
	}

	fromAccount, err := as.queries.GetAccountByID(ctx, fromAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}
	if !accountBelongsTo(fromAccount, holderID, holderType) {
		return ErrAccountNotOwned
	}

	tx, err := as.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := as.queries.WithTx(tx)

	firstID, secondID := fromAccountID, toAccountID
	if bytes.Compare(toAccountID[:], fromAccountID[:]) < 0 {
		firstID, secondID = toAccountID, fromAccountID
	}

	if _, err := qtx.GetAccountByIDForUpdate(ctx, firstID); err != nil {
		return err
	}
	if _, err := qtx.GetAccountByIDForUpdate(ctx, secondID); err != nil {
		return err
	}

	from, err := qtx.GetAccountByIDForUpdate(ctx, fromAccountID)
	if err != nil {
		return err
	}
	to, err := qtx.GetAccountByIDForUpdate(ctx, toAccountID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAccountNotFound
		}
		return err
	}

	if from.Status != pgstore.AccountStatusActive || to.Status != pgstore.AccountStatusActive {
		return ErrAccountNotActive
	}
	if from.Balance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	if err := qtx.UpdateAccountBalance(ctx, pgstore.UpdateAccountBalanceParams{
		ID:      fromAccountID,
		Balance: from.Balance.Sub(amount),
	}); err != nil {
		return err
	}
	if err := qtx.UpdateAccountBalance(ctx, pgstore.UpdateAccountBalanceParams{
		ID:      toAccountID,
		Balance: to.Balance.Add(amount),
	}); err != nil {
		return err
	}

	if _, err := qtx.CreateTransaction(ctx, pgstore.CreateTransactionParams{
		AccountID:             fromAccountID,
		Type:                  pgstore.TransactionTypeTransferOut,
		Amount:                amount,
		CounterpartyAccountID: pgtype.UUID{Bytes: toAccountID, Valid: true},
	}); err != nil {
		return err
	}
	if _, err := qtx.CreateTransaction(ctx, pgstore.CreateTransactionParams{
		AccountID:             toAccountID,
		Type:                  pgstore.TransactionTypeTransferIn,
		Amount:                amount,
		CounterpartyAccountID: pgtype.UUID{Bytes: fromAccountID, Valid: true},
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
