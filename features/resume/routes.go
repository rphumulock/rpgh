package resume

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	handlers := NewHandlers()

	router.Get("/resume", handlers.ResumePage)

	return nil
}
