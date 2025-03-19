package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"reminderai/repository"
)

type LogController struct {
	logRepo *repository.LogRepository
}

func NewLogController(logRepo *repository.LogRepository) *LogController {
	return &LogController{logRepo}
}

func (c *LogController) GetAll(w http.ResponseWriter, r *http.Request) {
	logs, err := c.logRepo.GetAll()
	if err != nil {
		log.Printf("Error retrieving logs: %v", err)
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
