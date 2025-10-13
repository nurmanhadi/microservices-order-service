package dto

import "time"

type ApiProductResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ApiResponse[T any] struct {
	Data T      `json:"data"`
	Path string `json:"path"`
}
