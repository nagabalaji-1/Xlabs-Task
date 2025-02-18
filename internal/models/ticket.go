package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials represents the login credentials.
type Credentials struct {
	Username string `json:"username"`
	//Password string `json:"password"`
}

// Train represents a train in the system.
type Train struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

// Ticket represents a train ticket.
type Ticket struct {
	ID      string       `json:"id"`
	TrainID string       `json:"train_id"`
	UserID  string       `json:"user_id"`
	SeatNo  int          `json:"seat_no"`
	Status  TicketStatus `json:"status"`
}

// TicketStatus represents the status of a ticket.
type TicketStatus int

const (
	// Enum values for TicketStatus
	Pending TicketStatus = iota
	Confirmed
	Cancelled
)

// String returns the string representation of the TicketStatus.
func (status TicketStatus) String() string {
	return [...]string{"Pending", "Confirmed", "Cancelled"}[status]
}
