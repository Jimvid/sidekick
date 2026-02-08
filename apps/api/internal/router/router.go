package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jimvid/sidekick/internal/config"
	"github.com/jimvid/sidekick/internal/database"
	"github.com/jimvid/sidekick/internal/habits"
	"github.com/jimvid/sidekick/internal/middleware"
)

func NewRouter(cfg *config.Config) *chi.Mux {

	r := chi.NewRouter()
	db := database.NewDynamoDB()

	// Habits
	habitStorage := habits.NewHabitStorage(db, cfg)
	habitService := habits.NewHabitService(habitStorage)
	habitHandler := habits.NewHabitHandler(habitService)

	// Cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Habits
	r.With(middleware.AuthMiddleware).Post("/habits", habitHandler.Create)
	r.With(middleware.AuthMiddleware).Get("/habits", habitHandler.GetAll)
	r.With(middleware.AuthMiddleware).Get("/habits/{habitId}", habitHandler.FindById)
	r.With(middleware.AuthMiddleware).Delete("/habits/{habitId}", habitHandler.Delete)
	r.With(middleware.AuthMiddleware).Put("/habits/{habitId}", habitHandler.Update)

	return r
}
