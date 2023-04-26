package store

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"

	"github.com/umalmyha/gonats/tickets-service/internal/model"
)

type TicketStore struct {
	db *sql.DB
}

func NewTicketStore(db *sql.DB) *TicketStore {
	return &TicketStore{db: db}
}

func (s *TicketStore) CreateTicket(ctx context.Context, t *model.Ticket) (*model.Ticket, error) {
	_, err := sq.Insert("tickets").
		Columns("id", "title", "description", "assigned_to", "created_at").
		Values(t.ID, t.Title, t.Description, t.AssignedTo, t.CreatedAt).
		RunWith(s.db).
		PlaceholderFormat(sq.Question).
		ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("TicketStore - CreateTicket - failed to insert new ticket: %w", err)
	}

	return &model.Ticket{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		AssignedTo:  t.AssignedTo,
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (s *TicketStore) GetAllTickets(ctx context.Context, query model.GetAllTicketsQuery) ([]*model.Ticket, error) {
	q := sq.Select("id", "title", "description", "assigned_to", "created_at").From("tickets")

	if query.Search != nil {
		search := *query.Search
		q = q.Where("(INSTR(title, ?) OR INSTR(title, ?) OR INSTR(assigned_to, ?))", search, search, search)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("TicketStore - GetAllTickets - failed to build get all tickets query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TicketStore - GetAllTickets - failed to query tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]*model.Ticket, 0)
	for rows.Next() {
		var t model.Ticket
		if err = rows.Scan(&t.ID, &t.Title, &t.Description, &t.AssignedTo, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("TicketStore - GetAllTickets - failed to scan ticket: %w", err)
		}
		tickets = append(tickets, &t)
	}

	return tickets, nil
}
