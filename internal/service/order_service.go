package service

import (
	"database/sql"
	"errors"
	"order-service/delivery/messaging/producer"
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
	AddOrder(request dto.OrderAddRequest) (*dto.PaymentResponse, error)
	GetAllOrder() ([]dto.OrderResponse, error)
	GetOrderByID(id string) (*dto.OrderResponse, error)
	GetAllOrderByUserID(userID string) ([]dto.OrderResponse, error)
	UpdateStatusByID(request dto.PaymentEventResponse) error
}
type orderService struct {
	logger          *logrus.Logger
	validation      *validator.Validate
	orderRepository repository.OrderRepository
	orderProducer   producer.OrderProducer
}

func NewOrderService(logger *logrus.Logger, validation *validator.Validate, orderRepository repository.OrderRepository, orderProducer producer.OrderProducer) OrderService {
	return &orderService{
		logger:          logger,
		validation:      validation,
		orderRepository: orderRepository,
		orderProducer:   orderProducer,
	}
}
func (s *orderService) AddOrder(request dto.OrderAddRequest) (*dto.PaymentResponse, error) {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("failed to validate request add order")
		return nil, err
	}
	orderID := uuid.NewString()
	var totalAmount int64
	items := make([]entity.Item, 0, len(request.Items))
	for _, x := range request.Items {
		product, err := api.GetProductByID(x.ProductID)
		if err != nil {
			s.logger.WithError(err).Error("failed to api get product by id")
			return nil, err
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
	paymentRequest := &dto.PaymentAddRequest{
		OrderID:       order.ID,
		TotalAmount:   order.TotalAmount,
		PaymentMethod: order.PaymentMethod,
	}
	response, err := api.PaymentCreateTransaction(*paymentRequest)
	if err != nil {
		s.logger.WithError(err).Error("failed to api payment create transaction")
		return nil, err
	}
	if err := s.orderRepository.Insert(*order); err != nil {
		s.logger.WithError(err).Error("failed to insert order")
		return nil, err
	}

	return response, nil
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
func (s *orderService) UpdateStatusByID(request dto.PaymentEventResponse) error {
	if err := s.validation.Struct(&request); err != nil {
		s.logger.WithError(err).Warn("failed to validate request update status by order id")
		return err
	}
	order, err := s.orderRepository.FindByID(request.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.WithError(err).Warn("order not found")
			return errors.New("order not found")
		}
		s.logger.WithError(err).Error("failed get order by id")
		return err
	}
	datas := make([]dto.OrderEventResponse, 0, len(order.Items))
	for _, x := range order.Items {
		datas = append(datas, dto.OrderEventResponse{
			OrderID:   x.OrderID,
			ProductID: x.ProductID,
			Quantity:  x.Quantity,
		})
	}
	err = s.orderRepository.UpdateStatusByID(order.ID, request.TransactionStatus)
	if err != nil {
		s.logger.WithError(err).Error("failed to update status by id")
		return err
	}
	go func() {
		err := s.orderProducer.PublishToOrderUpdated(datas)
		if err != nil {
			s.logger.WithError(err).Error("failed to publish to order updated")
			return
		}
	}()
	return nil
}
