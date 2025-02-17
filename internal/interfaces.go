package internal

import "go-ticket-app/internal/models"

// TicketStore defines the interface for ticket store operations.
type TicketStore interface {
	Create(ticket models.Ticket) (models.Ticket, error)
	Get(id string) (models.Ticket, bool)
	List() []models.Ticket
	Update(id string, updatedTicket models.Ticket) (models.Ticket, error)
	Delete(id string) error
}

// TicketQueue defines the interface for queue operations.
type TicketQueue interface {
	Enqueue(ticket models.Ticket)
	ProcessQueue()
	Get(id string) (models.Ticket, bool)
	List() []models.Ticket
	Update(id string, updatedTicket models.Ticket) (models.Ticket, error)
	Delete(id string) error
}
