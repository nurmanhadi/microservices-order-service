package dto

import (
	"order-service/pkg/enum"
)

type OrderAddRequest struct {
	UserID        string           `json:"user_id" validate:"required,max=36"`
	PaymentMethod enum.PAYMENT     `json:"payment_method" validate:"required,oneof=qris bcava"`
	Items         []ItemAddRequest `json:"items" validate:"required"`
}
