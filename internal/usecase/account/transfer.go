package account

import (
	"context"

	"github.com/jmarcosnsf/gobank/internal/validator"
)

type TransferReq struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

func (req TransferReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.FromAccountID), "from_account_id", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(req.ToAccountID), "to_account_id", "this field cannot be empty")
	eval.CheckField(req.Amount > 0, "amount", "must be greater than zero")

	return eval
}
