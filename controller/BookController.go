package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reminderai/model"
	"reminderai/repository"
	"strconv"
)

// BookController handles book-related HTTP endpoints
type BookController struct {
	bookRepo *repository.BookRepository
	logRepo  *repository.LogRepository
}

func NewBookController(bookRepo *repository.BookRepository, logRepo *repository.LogRepository) *BookController {
	return &BookController{bookRepo, logRepo}
}

func (c *BookController) Create(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Business logic can go here (validation, etc.)
	if book.Title == "" || book.Author == "" {
		http.Error(w, "Title and author are required", http.StatusBadRequest)
		return
	}

	if err := c.bookRepo.Create(&book); err != nil {
		log.Printf("Error creating book: %v", err)
		c.logRepo.Create("Failed to create book: "+err.Error(), "error")
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (c *BookController) GetAll(w http.ResponseWriter, r *http.Request) {
	books, err := c.bookRepo.GetAll()
	if err != nil {
		log.Printf("Error retrieving books: %v", err)
		c.logRepo.Create("Failed to retrieve books: "+err.Error(), "error")
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (c *BookController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	bookID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Business logic: validation
	if book.Title == "" || book.Author == "" {
		http.Error(w, "Title and author are required", http.StatusBadRequest)
		return
	}

	book.ID = bookID

	if err := c.bookRepo.Update(&book); err != nil {
		log.Printf("Error updating book: %v", err)
		c.logRepo.Create("Failed to update book: "+err.Error(), "error")
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}
