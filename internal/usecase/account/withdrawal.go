package account

import (
	"context"

	"github.com/jmarcosnsf/gobank/internal/validator"
)

type WithdrawalReq struct {
	Amount float64 `json:"amount"`
}

func (req WithdrawalReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(req.Amount > 0, "amount", "must be greater than zero")

	return eval
}