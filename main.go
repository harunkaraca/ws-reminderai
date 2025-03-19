package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"reminderai/controller"
	"reminderai/db"
	"reminderai/repository"
)

func main() {
	// Create a connection pool
	pool, err := db.NewPool()
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()
	// Print initial pool stats
	//db.PrintPoolStats(pool)

	// Optionally monitor pool stats every 30 seconds
	//db.MonitorPoolStats(pool, 1*time.Second)
	// Initialize the schema
	if err := db.InitSchema(pool); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Create repositories
	bookRepo := repository.NewBookRepository(pool)
	logRepo := repository.NewLogRepository(pool)

	// Log startup message
	if err := logRepo.Create("Server starting", "info"); err != nil {
		logRepo.Create("Failed to create startup log: "+err.Error(), "error")
	}

	// Create controllers
	bookController := controller.NewBookController(bookRepo, logRepo)
	logController := controller.NewLogController(logRepo)

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}).Methods("GET")

	router.HandleFunc("/books", bookController.Create).Methods("POST")
	router.HandleFunc("/books/{id}", bookController.Update).Methods("PUT")
	router.HandleFunc("/books", bookController.GetAll).Methods("GET")
	router.HandleFunc("/logs", logController.GetAll).Methods("GET")

	// Start the server
	log.Println("Starting server on :8080")
	port := ":" + os.Getenv("PORT")
	//port := ":8080"
	if err := http.ListenAndServe(port, router); err != nil {
		logRepo.Create("Server failed to start: "+err.Error(), "fatal")
		log.Fatalf("Failed to start server: %v", err)
	}
}

//Gin
//func main() {
//	router := gin.Default()
//	router.GET("/ping", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "pong",
//		})
//	})
//	router.Run() // listen and serve on 0.0.0.0:8080
//}
