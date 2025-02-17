package main

import (
    "net/http"

    "github.com/gorilla/mux"
    "go-ticket-app/internal/handlers"
)

func main() {
    router := mux.NewRouter()
    handlers.SetupRoutes(router)

    http.Handle("/", router)
    http.ListenAndServe(":8080", nil)
}