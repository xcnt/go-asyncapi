package amqp091go

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xcnt/go-asyncapi/run"
	runAmqp "github.com/xcnt/go-asyncapi/run/amqp"

	"github.com/rabbitmq/amqp091-go"
)

type SubscribeChannel struct {
	*amqp091.Channel
	// ConsumerTag uniquely identifies the consumer process. If empty, a unique tag is generated.
	ConsumerTag string
	// Additional arguments for the consumer. See ConsumeWithContext docs for details.
	ConsumeArgs amqp091.Table

	queueName string
	bindings  *runAmqp.ChannelBindings
}

func (s SubscribeChannel) Receive(ctx context.Context, cb func(envelope runAmqp.EnvelopeReader)) (err error) {
	// TODO: consumer tag in x- schema argument
	// Separate context is used to stop consumer process for a particular consumer tag on function exit.
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	deliveries, err := s.ConsumeWithContext(
		consumerCtx,
		s.queueName,
		s.ConsumerTag,
		s.bindings.SubscriberBindings.Ack,
		run.DerefOrZero(s.bindings.QueueConfiguration.Exclusive),
		false,
		false,
		s.ConsumeArgs,
	)
	if err != nil {
		return err
	}

	for delivery := range deliveries {
		evlp := NewEnvelopeIn(&delivery, bytes.NewReader(delivery.Body))
		cb(evlp)
		if s.bindings.SubscriberBindings.Ack {
			if e := s.Ack(delivery.DeliveryTag, false); e != nil {
				return fmt.Errorf("ack: %w", e)
			}
		}
	}
	return
}
