package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type CreateAssetRequest struct {
	WalletID uint    `json:"wallet_id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
}

func (r CreateAssetRequest) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(&r.WalletID, validation.Required),
		validation.Field(&r.Name, validation.Required),
	}

	return errors.Wrap(validation.ValidateStruct(&r, fields...), "deposit create validation error")
}
