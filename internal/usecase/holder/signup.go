package holder

import (
	"context"

	"github.com/jmarcosnsf/gobank/internal/validator"
)

type CreateIndividualReq struct {
	FullName      string  `json:"full_name"`
	Cpf           string  `json:"cpf"`
	DateOfBirth   string  `json:"date_of_birth"`
	Email         string  `json:"email"`
	Phone         string  `json:"phone"`
	Category      string  `json:"category"`
	MonthlyIncome float64 `json:"monthly_income"`
	Password      string  `json:"password"`
}

func (req CreateIndividualReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.FullName), "full_name", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.FullName, 255), "full_name", "must be at most 255 characters")

	eval.CheckField(validator.NotBlank(req.Cpf), "cpf", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.Cpf, 14), "cpf", "must be at most 14 characters")

	eval.CheckField(validator.NotBlank(req.DateOfBirth), "date_of_birth", "this field cannot be empty")
	eval.CheckField(validator.IsValidDate(req.DateOfBirth, "2006-01-02"), "date_of_birth", "must be a valid date in YYYY-MM-DD format")

	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "must be a valid email")
	eval.CheckField(validator.MaxChars(req.Email, 255), "email", "must be at most 255 characters")

	eval.CheckField(validator.NotBlank(req.Phone), "phone", "this field cannot be empty")
	eval.CheckField(validator.MaxChars(req.Phone, 20), "phone", "must be at most 20 characters")

	eval.CheckField(validator.MaxChars(req.Category, 50), "category", "must be at most 50 characters")

	eval.CheckField(validator.MinValue(req.MonthlyIncome, 0), "monthly_income", "must be greater than or equal to 0")

	eval.CheckField(validator.NotBlank(req.Password), "password", "this field cannot be empty")
	eval.CheckField(validator.MinChars(req.Password, 8), "password", "must be at least 8 characters")
	eval.CheckField(validator.MaxChars(req.Password, 72), "password", "must be at most 72 characters")

	return eval
}
