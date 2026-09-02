package portfolio

import (
	"rpgh/features/portfolio/components"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes serves the root listing and one route per directory on it, named
// the way the listing names it. The set comes from Dirs() rather than a list
// written out again here, so a directory added there is reachable without a
// second edit -- and there is no way to add one that the listing links to and
// nothing serves.
func SetupRoutes(router chi.Router) error {
	handlers := NewHandlers()

	router.Get("/", handlers.PortfolioPage(components.TabHome))
	for _, d := range components.Dirs() {
		router.Get(d.Href(), handlers.PortfolioPage(d.Key))
	}

	return nil
}
