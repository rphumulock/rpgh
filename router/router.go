package router

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"rpgh/config"
	blogFeature "rpgh/features/blog"
	hostFeature "rpgh/features/host"
	portfolioFeature "rpgh/features/portfolio"
	resumeFeature "rpgh/features/resume"
	"rpgh/web/resources"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(router chi.Router) (err error) {

	if config.Global.Environment == config.Dev {
		setupReload(router)
	}

	router.Handle("/static/*", resources.Handler())

	// Wired here rather than inside a feature because it answers for the whole
	// router: every directory is a URL now, so a wrong one is a page a visitor
	// reaches, and chi's bare 404 would be the one page on the site with none
	// of the site on it.
	router.NotFound(portfolioFeature.NotFoundHandler())

	if err := errors.Join(
		portfolioFeature.SetupRoutes(router),
		resumeFeature.SetupRoutes(router),
		hostFeature.SetupRoutes(router),
		blogFeature.SetupRoutes(router),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	return nil
}

func setupReload(router chi.Router) {
	reloadChan := make(chan struct{}, 1)
	var hotReloadOnce sync.Once

	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() { sse.ExecuteScript("window.location.reload()") }
		hotReloadOnce.Do(reload)
		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		}
	})

	router.Get("/hotreload", func(w http.ResponseWriter, r *http.Request) {
		select {
		case reloadChan <- struct{}{}:
		default:
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
