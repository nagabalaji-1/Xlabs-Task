package main

import (
	"bytes"
	"encoding/json"
	"go-ticket-app/internal/handlers"
	"go-ticket-app/internal/models"
	"go-ticket-app/internal/store"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestServer(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()

	// Mock the ticket store
	ticketStore := store.NewMapTicketStore()

	// Initialize handlers with the mocked store
	handlers.Init(ticketStore)
	handlers.SetupRoutes(router)

	// Create a test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Test CreateTicket endpoint with valid data
	ticket := models.Ticket{ID: "1", TrainID: "123", UserID: "456", SeatNo: 1, Status: models.Pending}
	ticketJSON, _ := json.Marshal(ticket)
	resp, err := http.Post(ts.URL+"/tickets", "application/json", bytes.NewBuffer(ticketJSON))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created; got %v", resp.Status)
	}

	// Test CreateTicket endpoint with invalid data
	resp, err = http.Post(ts.URL+"/tickets", "application/json", bytes.NewBuffer([]byte("invalid data")))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest; got %v", resp.Status)
	}

	// Test ListTickets endpoint when no tickets are available
	resp, err = http.Get(ts.URL + "/tickets")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status NotFound; got %v", resp.Status)
	}

	// Test ListTickets endpoint after creating a ticket
	resp, err = http.Get(ts.URL + "/tickets")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
}
