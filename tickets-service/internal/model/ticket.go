package model

import "time"

type Ticket struct {
	ID          string
	Title       string
	Description *string
	AssignedTo  string
	CreatedAt   time.Time
}

type GetAllTicketsQuery struct {
	Search *string
}
