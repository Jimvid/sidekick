package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jimvid/sidekick/internal/climbs"
	"github.com/jimvid/sidekick/internal/config"
	"github.com/jimvid/sidekick/internal/database"
	"github.com/jimvid/sidekick/internal/middleware"
)

func NewRouter(cfg *config.Config) *chi.Mux {

	r := chi.NewRouter()
	db := database.NewDynamoDB()

	// Climbs
	climbStorage := climbs.NewClimbStorage(db, cfg)
	climbService := climbs.NewClimbService(climbStorage)
	climbHandler := climbs.NewClimbHandler(climbService)

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

	// Climbs
	r.With(middleware.AuthMiddleware).Get("/climbs", climbHandler.GetAll)
	r.With(middleware.AuthMiddleware).Get("/climbs/{climbId}", climbHandler.FindById)
	r.With(middleware.AuthMiddleware).Post("/climbs", climbHandler.Create)
	r.With(middleware.AuthMiddleware).Delete("/climbs/{climbId}", climbHandler.Delete)
	r.With(middleware.AuthMiddleware).Post("/climbs/{climbId}", climbHandler.Update)

	return r
}
