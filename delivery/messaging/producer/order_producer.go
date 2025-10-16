package producer

import (
	"context"
	"encoding/json"
	"order-service/internal/entity"
	"order-service/pkg/env"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type OrderProducer interface {
	PublishToOrderCreated(order entity.Order) error
}
type orderProducer struct {
	ch *amqp091.Channel
}

func NewOrderProducer(ch *amqp091.Channel) OrderProducer {
	return &orderProducer{
		ch: ch,
	}
}
func (p *orderProducer) PublishToOrderCreated(order entity.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	for i := 0; i < 3; i++ {
		err = p.ch.PublishWithContext(
			ctx,
			env.CONF.Broker.Exchange.Order,
			env.CONF.Broker.RouteKey.OrderCreated,
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
		time.Sleep(500 * time.Microsecond)
	}
	return nil
}
