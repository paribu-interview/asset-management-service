package entity

import (
	"time"
)
import "gopkg.in/guregu/null.v3"

type Wallet struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt null.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
	Address   string    `json:"address"`
	Network   string    `json:"network"`
}
