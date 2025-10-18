package dto

import (
	"order-service/pkg/enum"
	"time"
)

type PaymentEventResponse struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
}
type PaymentAddRequest struct {
	OrderID       string       `json:"order_id"`
	TotalAmount   int64        `json:"total_amount"`
	PaymentMethod enum.PAYMENT `json:"payment_method"`
}
type PaymentResponse struct {
	ID                string    `json:"id"`
	OrderID           string    `json:"order_id"`
	GrossAmount       float64   `json:"gross_amount"`
	PaymentType       string    `json:"payment_type"`
	TransactionTime   time.Time `json:"transaction_time"`
	TransactionStatus string    `json:"transaction_status"`
	Currency          string    `json:"currency"`
	Bank              *string   `json:"bank,omitempty"`
	VaNumber          *string   `json:"va_number,omitempty"`
	Aquirer           *string   `json:"aquirer,omitempty"`
	QrString          *string   `json:"qr_string,omitempty"`
	ExpiryTime        time.Time `json:"expiry_time"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
