package repository

import (
	"fmt"
	"order-service/internal/entity"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Insert(order entity.Order) error
}
type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}
func (r *orderRepository) Insert(order entity.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO orders(id, user_id, total_amount, status, payment_method) VALUES($1, $2, $3, $4, $5)",
		order.ID, order.UserID, order.TotalAmount, order.Status, order.PaymentMethod)
	if err != nil {
		tx.Rollback()
		return err
	}
	query := "INSERT INTO items(order_id, product_id, quantity, price, subtotal) VALUES"
	for _, x := range order.Items {
		query += fmt.Sprintf("('%s', %d, %d, %d, %d),", x.OrderID, x.ProductID, x.Quantity, x.Price, x.Subtotal)
	}
	newQuery := strings.TrimSuffix(query, ",")
	fmt.Println(newQuery)
	_, err = tx.Exec(newQuery)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
