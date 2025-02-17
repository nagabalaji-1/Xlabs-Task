package queue

import (
	"go-ticket-app/internal"
	"go-ticket-app/internal/models"
	"sync"
)

type ticketQueue struct {
	queue []models.Ticket
	mu    sync.Mutex
	cond  *sync.Cond
	store internal.TicketStore
}

func NewTicketQueue(store internal.TicketStore) internal.TicketQueue {
	q := &ticketQueue{
		queue: make([]models.Ticket, 0),
		store: store,
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *ticketQueue) Enqueue(ticket models.Ticket) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, ticket)
	q.cond.Signal()
}

func (q *ticketQueue) Dequeue() models.Ticket {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.queue) == 0 {
		q.cond.Wait()
	}
	ticket := q.queue[0]
	q.queue = q.queue[1:]
	return ticket
}

func (q *ticketQueue) ProcessQueue() {
	for {
		ticket := q.Dequeue()
		q.store.Create(ticket)
	}
}

func (q *ticketQueue) Get(id string) (models.Ticket, bool) {
	return q.store.Get(id)
}

func (q *ticketQueue) List() []models.Ticket {
	return q.store.List()
}

func (q *ticketQueue) Update(id string, updatedTicket models.Ticket) (models.Ticket, error) {
	return q.store.Update(id, updatedTicket)
}

func (q *ticketQueue) Delete(id string) error {
	return q.store.Delete(id)
}
