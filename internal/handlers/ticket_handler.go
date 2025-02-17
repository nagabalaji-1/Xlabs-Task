package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "go-ticket-app/internal/models"
    "go-ticket-app/internal/store"
)

var ticketStore store.TicketStore

func Init(store store.TicketStore) {
    ticketStore = store
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
    var ticket models.Ticket
    if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createdTicket, err := ticketStore.Create(ticket)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTicket)
}

func GetTicket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    ticket, found := ticketStore.Get(id)
    if !found {
        http.Error(w, "Ticket not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(ticket)
}

func ListTickets(w http.ResponseWriter, r *http.Request) {
    allTickets := ticketStore.List()
    json.NewEncoder(w).Encode(allTickets)
}

func UpdateTicket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var updatedTicket models.Ticket
    if err := json.NewDecoder(r.Body).Decode(&updatedTicket); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    ticket, err := ticketStore.Update(id, updatedTicket)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(ticket)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    if err := ticketStore.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}