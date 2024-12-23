package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type CreateWithdrawRequest struct {
	WalletID uint    `json:"wallet_id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}

func (r CreateWithdrawRequest) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(&r.WalletID, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Amount, validation.Required),
	}

	return errors.Wrap(validation.ValidateStruct(&r, fields...), "withdraw create validation error")
}
