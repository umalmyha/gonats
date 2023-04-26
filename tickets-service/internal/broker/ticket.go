package broker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/velmie/broker"
	"github.com/velmie/broker/natsjs"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/umalmyha/gonats/tickets-service/internal/events"
)

const ticketCreatedSubject = "ticket.created"

type TicketEventsPublisher struct {
	pub *natsjs.Publisher
}

func NewTicketEventPublisher(pub *natsjs.Publisher) *TicketEventsPublisher {
	return &TicketEventsPublisher{pub: pub}
}

func (p *TicketEventsPublisher) RaiseTicketCreatedEvent(_ context.Context, evt events.TicketCreatedEvent) error {
	body, err := msgpack.Marshal(evt)
	if err != nil {
		return fmt.Errorf("TicketEventsPublisher - RaiseTicketCreatedEvent - failed to marshal ticket: %w", err)
	}

	err = p.pub.Publish(ticketCreatedSubject, &broker.Message{
		ID:     uuid.NewString(),
		Header: map[string]string{"Nats-Msg-Id": evt.ID},
		Body:   body,
	})
	if err != nil {
		return fmt.Errorf("TicketEventsPublisher - RaiseTicketCreatedEvent - failed to publish ticket created event: %w", err)
	}

	return nil
}
