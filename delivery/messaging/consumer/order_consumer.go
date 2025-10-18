package consumer

import (
	"encoding/json"
	"order-service/internal/dto"
	"order-service/internal/service"
	"order-service/pkg/env"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type OrderConsumer interface {
	ReceivedFromPaymentUpdated()
}
type orderConsumer struct {
	logger       *logrus.Logger
	Ch           *amqp091.Channel
	orderService service.OrderService
}

func NewOrderConsumer(logger *logrus.Logger, Ch *amqp091.Channel, orderService service.OrderService) OrderConsumer {
	return &orderConsumer{
		logger:       logger,
		Ch:           Ch,
		orderService: orderService,
	}
}
func (c *orderConsumer) ReceivedFromPaymentUpdated() {
	msgs, err := c.Ch.Consume(
		env.CONF.Broker.Queue.Order,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.WithError(err).Error("failed to consume queue")
	}
	go func() {
		for d := range msgs {
			request := new(dto.PaymentEventResponse)
			err := json.Unmarshal(d.Body, request)
			if err != nil {
				c.logger.WithError(err).Error("failed to unmarshal body")
			}
			err = c.orderService.UpdateStatusByID(*request)
			if err != nil {
				c.logger.WithError(err).Error("failed to update status by id")
			}
		}
	}()
}
