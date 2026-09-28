package main

import (
	"log"
	"net/http"

	"github.com/Ghahremanialireza/routinetask-backend/internal/database"
	"github.com/Ghahremanialireza/routinetask-backend/internal/router"
)

func main() {
	db, err := database.New("routinetask.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	taskRepo := database.NewTaskRepository(db)

	r := router.New(taskRepo)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
