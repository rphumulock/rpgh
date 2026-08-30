package host

import (
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	handlers := NewHandlers()

	router.Get("/host/events", handlers.Events)

	return nil
}
