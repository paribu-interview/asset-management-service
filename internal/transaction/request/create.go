package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"time"
)

type ScheduleTransactionRequest struct {
	SourceWalletID      uint      `json:"source_wallet_id"`
	DestinationWalletID uint      `json:"destination_wallet_id"`
	AssetName           string    `json:"asset_name"`
	Amount              float64   `json:"amount"`
	ScheduledAt         time.Time `json:"scheduled_at"`
}

func (r ScheduleTransactionRequest) Validate() error {
	fields := []*validation.FieldRules{
		validation.Field(&r.SourceWalletID, validation.Required),
		validation.Field(&r.DestinationWalletID, validation.Required),
		validation.Field(&r.AssetName, validation.Required),
		validation.Field(&r.Amount, validation.Required),
		validation.Field(&r.ScheduledAt, validation.Required),
	}

	return errors.Wrap(validation.ValidateStruct(&r, fields...), "schedule transaction create validation error")
}
