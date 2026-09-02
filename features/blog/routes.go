package blog

import (
	"fmt"

	"rpgh/features/blog/content"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes wires the post pages, and refuses to start if posts/ does not
// parse. A malformed post is a mistake made once, at the keyboard; finding out
// at boot with the filename in the error beats finding out from a visitor.
func SetupRoutes(router chi.Router) error {
	if err := content.Err(); err != nil {
		return fmt.Errorf("reading the blog: %w", err)
	}

	handlers := NewHandlers()

	router.Get("/blog/{slug}", handlers.PostPage)

	return nil
}
