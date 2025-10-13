package dto

import "time"

type ItemAddRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	Quantity  int   `json:"quantity" validate:"required"`
}
type ItemResponse struct {
	ID        int64          `json:"id"`
	OrderID   string         `json:"order_id"`
	ProductID int64          `json:"product_id"`
	Quantity  int            `json:"quantity"`
	Price     int64          `json:"price"`
	Subtotal  int64          `json:"subtotal"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Order     *OrderResponse `json:"orders"`
}
