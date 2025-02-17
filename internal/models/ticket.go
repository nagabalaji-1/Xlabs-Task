package models

// Ticket struct represents a ticket with an ID, title, and description.
type Ticket struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}