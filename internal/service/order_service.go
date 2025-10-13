package service

import (
	"database/sql"
	"order-service/internal/dto"
	"order-service/internal/entity"
	"order-service/internal/repository"
	"order-service/pkg/api"
	"order-service/pkg/enum"
	"order-service/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderService interface {
	AddOrder(request dto.OrderAddRequest) error
	GetAllOrder() ([]dto.OrderResponse, error)
	GetOrderByID(id string) (*dto.OrderResponse, error)
	GetAllOrderByUserID(userID string) ([]dto.OrderResponse, error)
}
type orderService struct {
	logger          *logrus.Logger
	validation      *validator.Validate
	orderRepository repository.OrderRepository
	productAPI      api.ProductAPI
}

func NewOrderService(logger *logrus.Logger, validation *validator.Validate, orderRepository repository.OrderRepository, productAPI api.ProductAPI) OrderService {
	return &orderService{
		logger:          logger,
		validation:      validation,
		orderRepository: orderRepository,
		productAPI:      productAPI,
	}
}
func (s *orderService) AddOrder(request dto.OrderAddRequest) error {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("failed to validate request add order")
		return err
	}
	orderID := uuid.NewString()
	var totalAmount int64
	items := make([]entity.Item, 0, len(request.Items))
	for _, x := range request.Items {
		product, err := s.productAPI.GetProductByID(x.ProductID)
		if err != nil {
			s.logger.WithError(err).Error("failed to api get product by id")
			return err
		}
		items = append(items, entity.Item{
			OrderID:   orderID,
			ProductID: product.ID,
			Quantity:  x.Quantity,
			Price:     product.Price,
			Subtotal:  int64(x.Quantity) * product.Price,
		})
		totalAmount += int64(x.Quantity) * product.Price
	}
	order := &entity.Order{
		ID:            orderID,
		UserID:        request.UserID,
		TotalAmount:   totalAmount,
		Status:        enum.STATUS_PENDING,
		PaymentMethod: request.PaymentMethod,
		Items:         items,
	}
	if err := s.orderRepository.Insert(*order); err != nil {
		s.logger.WithError(err).Error("failed to insert order")
		return err
	}
	return nil
}
func (s *orderService) GetAllOrder() ([]dto.OrderResponse, error) {
	orders, err := s.orderRepository.FindAll()
	if err != nil {
		s.logger.WithError(err).Error("failed to get all order")
		return nil, err
	}
	resp := make([]dto.OrderResponse, 0, len(orders))
	for _, x := range orders {
		resp = append(resp, dto.OrderResponse{
			ID:            x.ID,
			UserID:        x.UserID,
			TotalAmount:   x.TotalAmount,
			Status:        x.Status,
			PaymentMethod: x.PaymentMethod,
			CreatedAt:     x.CreatedAt,
			UpdatedAt:     x.CreatedAt,
		})
	}
	return resp, nil
}
func (s *orderService) GetOrderByID(id string) (*dto.OrderResponse, error) {
	order, err := s.orderRepository.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.WithError(err).Warn("order not found")
			return nil, response.Except(404, "order not found")
		}
		s.logger.WithError(err).Error("failed get order by id")
		return nil, err
	}
	items := make([]dto.ItemResponse, 0, len(order.Items))
	for _, x := range order.Items {
		items = append(items, dto.ItemResponse{
			ID:        x.ID,
			OrderID:   x.OrderID,
			ProductID: x.ProductID,
			Quantity:  x.Quantity,
			Price:     x.Price,
			Subtotal:  x.Subtotal,
			CreatedAt: x.CreatedAt,
			UpdatedAt: x.UpdatedAt,
		})
	}
	resp := &dto.OrderResponse{
		ID:            order.ID,
		UserID:        order.UserID,
		TotalAmount:   order.TotalAmount,
		Status:        order.Status,
		PaymentMethod: order.PaymentMethod,
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
		Items:         items,
	}
	return resp, nil
}
func (s *orderService) GetAllOrderByUserID(userID string) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepository.FindAllByUserID(userID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get all order by user id")
		return nil, err
	}
	resp := make([]dto.OrderResponse, 0, len(orders))
	for _, x := range orders {
		resp = append(resp, dto.OrderResponse{
			ID:            x.ID,
			UserID:        x.UserID,
			TotalAmount:   x.TotalAmount,
			Status:        x.Status,
			PaymentMethod: x.PaymentMethod,
			CreatedAt:     x.CreatedAt,
			UpdatedAt:     x.CreatedAt,
		})
	}
	return resp, nil
}
