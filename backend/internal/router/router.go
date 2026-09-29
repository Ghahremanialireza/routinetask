package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Ghahremanialireza/routinetask-backend/internal/database"
	"github.com/Ghahremanialireza/routinetask-backend/internal/handlers"
)

func New(taskRepo *database.TaskRepository, allowedOrigin string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	}))

	r.Get("/health", handlers.HealthHandler)

	taskHandler := handlers.NewTaskHandler(taskRepo)
	r.Post("/tasks", taskHandler.CreateTask)
	r.Get("/tasks", taskHandler.GetTasks)
	r.Patch("/tasks/{id}", taskHandler.UpdateTaskCompleted)
	r.Delete("/tasks/{id}", taskHandler.DeleteTask)

	return r
}
