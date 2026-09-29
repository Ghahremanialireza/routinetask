package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ghahremanialireza/routinetask-backend/internal/database"
	"github.com/Ghahremanialireza/routinetask-backend/internal/router"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	dbPath := getEnv("DB_PATH", "routinetask.db")
	port := getEnv("PORT", "8080")
	allowedOrigin := getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:3000")

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	taskRepo := database.NewTaskRepository(db)
	r := router.New(taskRepo, allowedOrigin)

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
