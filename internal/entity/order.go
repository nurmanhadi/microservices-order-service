package entity

import (
	"order-service/pkg/enum"
	"time"
)

type Order struct {
	ID            string       `db:"id"`
	UserID        string       `db:"user_id"`
	TotalAmount   int64        `db:"total_amount"`
	Status        enum.STATUS  `db:"status"`
	PaymentMethod enum.PAYMENT `db:"payment_method"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	Items         []Item       `db:"items"`
}
