package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"go-ticket-app/internal"
	"go-ticket-app/internal/errors"
	"go-ticket-app/internal/models"
	"go-ticket-app/internal/queue"

	"github.com/gorilla/mux"
)

var (
	users       = make(map[string]models.User)
	trains      = make(map[string]models.Train)
	tickets     = make(map[string]models.Ticket)
	userMu      sync.Mutex
	trainMu     sync.Mutex
	ticketMu    sync.Mutex
	ticketQueue internal.TicketQueue
	logger      *log.Logger
)

func SetLogger(l *log.Logger) {
	logger = l
}

func Init(store internal.TicketStore) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file")
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ticketQueue = queue.NewTicketQueue(store)
	go ticketQueue.ProcessQueue()
}

func CreateTicket(w http.ResponseWriter, r *http.Request) {
	logger.Println("CreateTicket endpoint called")
	var ticket models.Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		logger.Printf("Error decoding ticket: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Printf("Decoded ticket: %+v", ticket)
	ticket.Status = models.Pending
	ticketQueue.Enqueue(ticket)
	w.WriteHeader(http.StatusAccepted)
	logger.Println("Ticket enqueued successfully")
}

func GetTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ticket, found := ticketQueue.Get(id)
	if !found {
		logger.Printf("Ticket not found: %v", id)
		http.Error(w, (&errors.NotFoundError{Resource: "Ticket", ID: id}).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedTicket models.Ticket
	if err := json.NewDecoder(r.Body).Decode(&updatedTicket); err != nil {
		logger.Printf("Error decoding updated ticket: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := ticketQueue.Update(id, updatedTicket)
	if err != nil {
		logger.Printf("Error updating ticket: %v", err)
		http.Error(w, (&errors.NotFoundError{Resource: "Ticket", ID: id}).Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := ticketQueue.Delete(id); err != nil {
		logger.Printf("Error deleting ticket: %v", err)
		http.Error(w, (&errors.NotFoundError{Resource: "Ticket", ID: id}).Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Register handles user registration.
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Printf("Error decoding user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	if _, exists := users[user.ID]; exists {
		logger.Printf("User already exists: %v", user.ID)
		http.Error(w, (&errors.ConflictError{Resource: "User", ID: user.ID}).Error(), http.StatusConflict)
		return
	}

	users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user login.
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		logger.Printf("Error decoding credentials: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	for _, user := range users {
		if user.Username == credentials.Username {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
			return
		}
	}

	http.Error(w, "Invalid user", http.StatusUnauthorized)
}

// UpdateUser handles updating a user's information.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		logger.Printf("Error decoding updated user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userMu.Lock()
	defer userMu.Unlock()

	if _, exists := users[id]; !exists {
		logger.Printf("User not found: %v", id)
		http.Error(w, (&errors.NotFoundError{Resource: "User", ID: id}).Error(), http.StatusNotFound)
		return
	}

	users[id] = updatedUser
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully", "id": id})
}

// DeleteUser handles deleting a user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userMu.Lock()
	defer userMu.Unlock()

	if _, exists := users[id]; !exists {
		logger.Printf("User not found: %v", id)
		http.Error(w, (&errors.NotFoundError{Resource: "User", ID: id}).Error(), http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}

// ListUsers handles listing all users.
func ListUsers(w http.ResponseWriter, r *http.Request) {
	userMu.Lock()
	defer userMu.Unlock()

	if len(users) > 0 {
		fmt.Fprintln(w, "User Details are :")
		log.Println("User Details are :")
		w.WriteHeader(http.StatusAccepted)
		for _, user := range users {
			fmt.Fprintf(w, " Id: %v, Username: %v \n", user.ID, user.Username)
			log.Printf(" Id: %v, Username: %v \n", user.ID, user.Username)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.User{})
	}
}

// CreateTrain handles creating a new train.
func CreateTrain(w http.ResponseWriter, r *http.Request) {
	var train models.Train
	if err := json.NewDecoder(r.Body).Decode(&train); err != nil {
		logger.Printf("Error decoding train: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trainMu.Lock()
	defer trainMu.Unlock()

	if _, exists := trains[train.ID]; exists {
		logger.Printf("Train already exists: %v", train.ID)
		http.Error(w, (&errors.ConflictError{Resource: "Train", ID: train.ID}).Error(), http.StatusConflict)
		return
	}

	trains[train.ID] = train
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Train created successfully"})
}

// ListTrains handles listing all trains.
func ListTrains(w http.ResponseWriter, r *http.Request) {
	trainMu.Lock()
	defer trainMu.Unlock()

	if len(trains) > 0 {
		fmt.Fprintln(w, "Train Details are :")
		log.Println("Train Details are :")
		w.WriteHeader(http.StatusAccepted)
		for _, train := range trains {
			fmt.Fprintf(w, " Id: %v, Name: %v, Capacity: %v \n", train.ID, train.Name, train.Capacity)
			log.Printf(" Id: %v, Name: %v, Capacity: %v \n", train.ID, train.Name, train.Capacity)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Train{})
	}
}

// BookTicket handles booking a ticket.
func BookTicket(w http.ResponseWriter, r *http.Request) {
	var ticket models.Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		logger.Printf("Error decoding ticket: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trainMu.Lock()
	train, exists := trains[ticket.TrainID]
	trainMu.Unlock()

	if !exists {
		logger.Printf("Train not found: %v", ticket.TrainID)
		http.Error(w, (&errors.NotFoundError{Resource: "Train", ID: ticket.TrainID}).Error(), http.StatusNotFound)
		return
	}

	ticketMu.Lock()
	defer ticketMu.Unlock()

	if len(tickets) >= train.Capacity {
		logger.Printf("Train is fully booked: %v", ticket.TrainID)
		http.Error(w, "Train is fully booked", http.StatusConflict)
		return
	}

	ticket.ID = fmt.Sprintf("Ticket%d", len(tickets)+1)
	tickets[ticket.ID] = ticket
	ticketQueue.Enqueue(ticket)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket booked successfully", "ticket_id": ticket.ID})
}

// ListTickets handles listing all tickets.
func ListTickets(w http.ResponseWriter, r *http.Request) {
	ticketMu.Lock()
	defer ticketMu.Unlock()
	logger.Println("ListTickets endpoint called")
	if len(tickets) > 0 {
		logger.Println("Ticket Details are :")
		fmt.Fprintln(w, "Ticket Details are :")
		log.Println("Ticket Details are :")
		w.WriteHeader(http.StatusAccepted)
		for _, ticket := range tickets {
			fmt.Fprintf(w, " Id: %v, TrainID: %v, UserID: %v, SeatNo: %v \n", ticket.ID, ticket.TrainID, ticket.UserID, ticket.SeatNo)
			log.Printf(" Id: %v, TrainID: %v, UserID: %v, SeatNo: %v \n", ticket.ID, ticket.TrainID, ticket.UserID, ticket.SeatNo)
		}
	} else {
		logger.Println("No tickets found")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Ticket{})
	}
}
