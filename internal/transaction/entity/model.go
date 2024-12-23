package entity

import "time"

type Filters struct {
	ID                  []uint
	SourceWalletID      []uint
	DestinationWalletID []uint
	Status              []string
	ScheduledStart      time.Time
	ScheduledEnd        time.Time
}
