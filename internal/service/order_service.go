package service

import (
	"order-service/internal/dto"
	"order-service/internal/entity"
	"order-service/internal/repository"
	"order-service/pkg/api"
	"order-service/pkg/enum"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OrderService interface {
	AddOrder(request dto.OrderAddRequest) error
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
