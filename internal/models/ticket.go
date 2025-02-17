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
	ID      string `json:"id"`
	TrainID string `json:"train_id"`
	UserID  string `json:"user_id"`
	SeatNo  int    `json:"seat_no"`
}
