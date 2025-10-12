package entity

import "time"

type Item struct {
	ID        int64     `db:"id"`
	OrderID   string    `db:"order_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	Price     int64     `db:"price"`
	Subtotal  int64     `db:"subtotal"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Order     *Order    `db:"orders"`
}
