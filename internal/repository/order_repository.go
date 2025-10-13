package repository

import (
	"fmt"
	"order-service/internal/entity"
	"strings"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Insert(order entity.Order) error
	FindAll() ([]entity.Order, error)
	FindByID(id string) (*entity.Order, error)
	FindAllByUserID(userID string) ([]entity.Order, error)
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
func (r *orderRepository) FindAll() ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Select(&orders, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (r *orderRepository) FindByID(id string) (*entity.Order, error) {
	order := new(entity.Order)
	err := r.db.Get(order, "SELECT * FROM orders WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	var items []entity.Item
	err = r.db.Select(&items, "SELECT * FROM items WHERE order_id = $1", id)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}
func (r *orderRepository) FindAllByUserID(userID string) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Select(&orders, "SELECT * FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
