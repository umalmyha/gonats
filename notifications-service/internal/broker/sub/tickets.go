package sub

import (
	"context"
	"fmt"
	"github.com/umalmyha/gonats/notifications-service/internal/events"
	"github.com/velmie/broker"
	"github.com/vmihailenco/msgpack/v5"
)

type TicketEventsSubscriber struct {
	sub broker.Subscriber
}

func NewTicketEventsSubscriber(sub broker.Subscriber) *TicketEventsSubscriber {
	return &TicketEventsSubscriber{sub: sub}
}

func (s *TicketEventsSubscriber) Listen(ctx context.Context, topic string) (err error) {
	sub, err := s.sub.Subscribe(topic, s.handler, broker.DisableAutoAck())
	if err != nil {
		return err
	}

	<-ctx.Done()
	err = sub.Unsubscribe()
	<-sub.Done()

	return
}

func (s *TicketEventsSubscriber) handler(evt broker.Event) error {
	m := evt.Message()

	var tc events.TicketCreatedEvent
	if err := msgpack.Unmarshal(m.Body, &tc); err != nil {
		return fmt.Errorf("TicketEventsSubscriber - handler - failed to unmarshal message body: %w", err)
	}

	fmt.Printf("Ticket created event received %v\n", tc)

	if err := evt.Ack(); err != nil {
		return fmt.Errorf("TicketEventsSubscriber - handler - failed to ack message: %w", err)
	}

	return nil
}
