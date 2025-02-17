package handlers

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures the HTTP routes for the ticket operations.
func SetupRoutes(router *mux.Router) {
	router.HandleFunc("/ticket", CreateTicket).Methods("POST")
	router.HandleFunc("/ticket/{id}", GetTicket).Methods("GET")
	router.HandleFunc("/tickets", ListTickets).Methods("GET")
	router.HandleFunc("/ticket/{id}", UpdateTicket).Methods("PUT")
	router.HandleFunc("/ticket/{id}", DeleteTicket).Methods("DELETE")

	// Add routes for Login and Register
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// Add routes for UpdateUser, DeleteUser, and ListUsers
	router.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/users", ListUsers).Methods("GET")

	// Add routes for Train operations
	router.HandleFunc("/train", CreateTrain).Methods("POST")
	router.HandleFunc("/trains", ListTrains).Methods("GET")

	// Add routes for Ticket operations
	router.HandleFunc("/book", BookTicket).Methods("POST")
	router.HandleFunc("/tickets", ListTickets).Methods("GET")
}
