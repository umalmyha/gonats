package events

type TicketCreatedEvent struct {
	ID         string
	Title      string
	AssignedTo string
}
