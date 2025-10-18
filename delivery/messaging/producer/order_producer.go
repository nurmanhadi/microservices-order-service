package producer

import (
	"context"
	"encoding/json"
	"order-service/internal/dto"
	"order-service/pkg/env"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type OrderProducer interface {
	PublishToOrderUpdated(datas []dto.OrderEventResponse) error
}
type orderProducer struct {
	ch *amqp091.Channel
}

func NewOrderProducer(ch *amqp091.Channel) OrderProducer {
	return &orderProducer{
		ch: ch,
	}
}
func (p *orderProducer) PublishToOrderUpdated(datas []dto.OrderEventResponse) error {
	body, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	err = p.ch.PublishWithContext(
		ctx,
		env.CONF.Broker.Exchange.Order,
		env.CONF.Broker.RouteKey.OrderUpdated,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
