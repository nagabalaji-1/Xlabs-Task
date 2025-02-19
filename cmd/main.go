package main

import (
	"go-ticket-app/internal/handlers"
	"go-ticket-app/internal/store"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	handlers.SetLogger(logger)
}

func main() {
	router := mux.NewRouter()
	ticketStore := store.NewMapTicketStore()
	handlers.Init(ticketStore)
	handlers.SetupRoutes(router)
	logger.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatalf("Could not start server: %s\n", err.Error())
	}
}
