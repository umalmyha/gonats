package ticketserver

import (
	"context"
	"github.com/umalmyha/gonats/tickets-service/internal/broker"
	"github.com/umalmyha/gonats/tickets-service/internal/events"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/umalmyha/gonats/tickets-service/internal/model"
	"github.com/umalmyha/gonats/tickets-service/internal/store"
	pb "github.com/umalmyha/gonats/tickets-service/rpc/ticket"
)

type Server struct {
	ticketStore     *store.TicketStore
	ticketPublisher *broker.TicketEventsPublisher
}

func NewServer(
	ticketStore *store.TicketStore,
	ticketPublisher *broker.TicketEventsPublisher,
) *Server {
	return &Server{
		ticketStore:     ticketStore,
		ticketPublisher: ticketPublisher,
	}
}

func (s *Server) CreateTicket(ctx context.Context, req *pb.CreateTicketRequest) (*pb.CreateTicketResponse, error) {
	nt, err := s.ticketStore.CreateTicket(ctx, &model.Ticket{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	err = s.ticketPublisher.RaiseTicketCreatedEvent(ctx, events.TicketCreatedEvent{
		ID:         nt.ID,
		Title:      nt.Title,
		AssignedTo: nt.AssignedTo,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateTicketResponse{Ticket: s.ticketToProto(nt)}, nil
}

func (s *Server) GetAllTickets(ctx context.Context, req *pb.GetAllTicketsRequest) (*pb.GetAllTicketsResponse, error) {
	tickets, err := s.ticketStore.GetAllTickets(ctx, model.GetAllTicketsQuery{
		Search: req.Search,
	})
	if err != nil {
		return nil, err
	}

	res := &pb.GetAllTicketsResponse{Tickets: make([]*pb.Ticket, 0)}
	for _, t := range tickets {
		res.Tickets = append(res.Tickets, s.ticketToProto(t))
	}

	return res, nil
}

func (s *Server) ticketToProto(t *model.Ticket) *pb.Ticket {
	return &pb.Ticket{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedAt:   timestamppb.New(t.CreatedAt),
	}
}
