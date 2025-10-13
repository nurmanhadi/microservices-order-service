package dto

import (
	"order-service/pkg/enum"
	"time"
)

type OrderAddRequest struct {
	UserID        string           `json:"user_id" validate:"required,max=36"`
	PaymentMethod enum.PAYMENT     `json:"payment_method" validate:"required,oneof=qris bcava"`
	Items         []ItemAddRequest `json:"items" validate:"required"`
}
type OrderResponse struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	TotalAmount   int64          `json:"total_amount"`
	Status        enum.STATUS    `json:"status"`
	PaymentMethod enum.PAYMENT   `json:"payment_method"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Items         []ItemResponse `json:"items"`
}
