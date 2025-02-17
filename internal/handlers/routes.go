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
}