package main

import (
	"net/http"

	"go-ticket-app/internal/handlers"
	"go-ticket-app/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	ticketStore := store.NewMapTicketStore()
	handlers.Init(ticketStore)
	handlers.SetupRoutes(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
