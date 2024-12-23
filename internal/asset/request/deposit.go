package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type CreateDepositRequest struct {
	WalletID uint    `json:"wallet_id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}

func (r CreateDepositRequest) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(&r.WalletID, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Amount, validation.Required),
	}

	return errors.Wrap(validation.ValidateStruct(&r, fields...), "deposit create validation error")
}
