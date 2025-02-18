package store

import (
	"fmt"
	"go-ticket-app/internal"
	"go-ticket-app/internal/models"
	"sync"
)

// MapTicketStore implements TicketStore using a map with a mutex for concurrency.
type MapTicketStore struct {
	tickets map[string]models.Ticket
	mu      sync.RWMutex
}

// NewMapTicketStore initializes a new MapTicketStore.
func NewMapTicketStore() internal.TicketStore {
	return &MapTicketStore{tickets: make(map[string]models.Ticket)}
}

// Create adds a new ticket to the store.
func (store *MapTicketStore) Create(ticket models.Ticket) (models.Ticket, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	ticketID := fmt.Sprintf("Ticket%d", len(store.tickets)+1)
	ticket.ID = ticketID
	ticket.Status = models.Pending
	store.tickets[ticketID] = ticket
	return ticket, nil
}

// Get retrieves a ticket by its ID.
func (store *MapTicketStore) Get(id string) (models.Ticket, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	ticket, found := store.tickets[id]
	return ticket, found
}

// List returns all tickets in the store.
func (store *MapTicketStore) List() []models.Ticket {
	store.mu.RLock()
	defer store.mu.RUnlock()

	allTickets := make([]models.Ticket, 0, len(store.tickets))
	for _, ticket := range store.tickets {
		allTickets = append(allTickets, ticket)
	}
	return allTickets
}

// Update modifies an existing ticket by its ID.
func (store *MapTicketStore) Update(id string, updatedTicket models.Ticket) (models.Ticket, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	ticket, found := store.tickets[id]
	if !found {
		return models.Ticket{}, fmt.Errorf("ticket not found")
	}
	ticket.ID = updatedTicket.ID
	ticket.SeatNo = updatedTicket.SeatNo
	ticket.TrainID = updatedTicket.TrainID
	ticket.UserID = updatedTicket.UserID
	ticket.Status = updatedTicket.Status
	store.tickets[id] = ticket
	return ticket, nil
}

// Delete removes a ticket from the store by its ID.
func (store *MapTicketStore) Delete(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, found := store.tickets[id]; !found {
		return fmt.Errorf("ticket not found")
	}
	delete(store.tickets, id)
	return nil
}
