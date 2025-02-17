// filepath: /go-ticket-app/go-ticket-app/internal/store/ticket_store.go
package store

import "go-ticket-app/internal/models"

// TicketStore interface defines the methods for ticket storage.
type TicketStore interface {
    Create(ticket models.Ticket) (models.Ticket, error)
    Get(id string) (models.Ticket, bool)
    List() []models.Ticket
    Update(id string, updatedTicket models.Ticket) (models.Ticket, error)
    Delete(id string) error
}