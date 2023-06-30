package models

import "time"

type Trx struct {
	ID              string
	SenderAccount   string
	ReceiverAccount string
	Type            string
	Amount          int
	CreatedAt       time.Time
}
