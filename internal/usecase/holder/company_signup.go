package holder

import (
	"context"

	"github.com/jmarcosnsf/gobank/internal/validator"
)

type CreateCompanyReq struct {
	TradeName      string  `json:"trade_name"`
	Cnpj           string  `json:"cnpj"`
	FoundedAt      string  `json:"founded_at"`
	CorporateEmail string  `json:"corporate_email"`
	Phone          string  `json:"phone"`
	Category       string  `json:"category"`
	AnnualRevenue  float64 `json:"annual_revenue"`
	Password       string  `json:"password"`
}

func (req CreateCompanyReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.TradeName), "trade_name", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.TradeName, 255), "trade_name", "must be at most 255 characters")

	eval.CheckField(validator.NotBlank(req.Cnpj), "cnpj", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.Cnpj, 18), "cnpj", "must be at most 18 characters")

	eval.CheckField(validator.NotBlank(req.FoundedAt), "founded_at", "this field cannot be empty")
	eval.CheckField(validator.IsValidDate(req.FoundedAt, "2006-01-02"), "founded_at", "must be a valid date in YYYY-MM-DD format")

	eval.CheckField(validator.NotBlank(req.CorporateEmail), "corporate_email", "this field cannot be empty")
	eval.CheckField(validator.Matches(req.CorporateEmail, validator.EmailRX), "corporate_email", "must be a valid email")
	eval.CheckField(validator.MaxChars(req.CorporateEmail, 255), "corporate_email", "must be at most 255 characters")

	eval.CheckField(validator.NotBlank(req.Phone), "phone", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.Phone, 20), "phone", "must be at most 20 characters")

	eval.CheckField(validator.MaxChars(req.Category, 50), "category", "must be at most 50 characters")

	eval.CheckField(validator.MinValue(req.AnnualRevenue, 0), "annual_revenue", "must be greater than or equal to 0")

	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.Password, 8), "password", "must be at least 8 characters")
	eval.CheckField(validator.MaxChars(req.Password, 72), "password", "must be at most 72 characters")

	return eval
}
