package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

type GetTransactionsParams struct {
	ID                  []uint   `json:"id" schema:"id"`
	SourceWalletID      []uint   `json:"source_wallet_id" schema:"source_wallet_id"`
	DestinationWalletID []uint   `json:"destination_wallet_id" schema:"destination_wallet_id"`
	Status              []string `json:"status" schema:"status"`
}

func (r GetTransactionsParams) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(&r.Status, validation.By(func(value interface{}) error {
			if value == nil || len(value.([]string)) == 0 {
				return nil
			}

			for _, v := range value.([]string) {
				if v != "pending" && v != "completed" && v != "cancelled" && v != "failed" {
					return errors.New("invalid status")
				}
			}
			return nil
		})),
	}

	return errors.Wrap(validation.ValidateStruct(&r, fields...), "schedule transaction create validation error")
}
