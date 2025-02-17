package handlers

import (
	"encoding/json"
	"net/http"

	"go-ticket-app/internal"
	"go-ticket-app/internal/models"
	"go-ticket-app/internal/queue"

	"github.com/gorilla/mux"
)

var ticketQueue internal.TicketQueue

func Init(store internal.TicketStore) {
	ticketQueue = queue.NewTicketQueue(store)
	go ticketQueue.ProcessQueue()
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticketQueue.Enqueue(ticket)
	w.WriteHeader(http.StatusAccepted)
}

func GetTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ticket, found := ticketQueue.Get(id)
	if !found {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func ListTickets(w http.ResponseWriter, r *http.Request) {
	allTickets := ticketQueue.List()
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

	ticket, err := ticketQueue.Update(id, updatedTicket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := ticketQueue.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
