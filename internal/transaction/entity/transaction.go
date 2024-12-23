package entity

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type Transaction struct {
	ID                  uint              `json:"id"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           null.Time         `json:"updated_at"`
	SourceWalletID      uint              `json:"source_wallet_id"`
	DestinationWalletID uint              `json:"destination_wallet_id"`
	AssetName           string            `json:"asset_name"`
	Amount              float64           `json:"amount"`
	Status              TransactionStatus `json:"status"`
	ScheduledAt         time.Time         `json:"scheduled_at"`
}

func (Transaction) TableName() string {
	return "scheduled_transactions"
}

type TransactionStatus string

const (
	TransactionCompleted TransactionStatus = "completed"
	TransactionFailed    TransactionStatus = "failed"
	TransactionPending   TransactionStatus = "pending"
	TransactionCancelled TransactionStatus = "cancelled"
)
