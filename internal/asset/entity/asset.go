package entity

import (
	"time"
)
import "gopkg.in/guregu/null.v3"

type Asset struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	WalletID  uint      `json:"wallet_id"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
}
