package portfolio

import (
	"net/http"

	"rpgh/features/portfolio/components"
	"rpgh/features/portfolio/pages"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// PortfolioPage serves one directory of the site. Which one is decided at
// wiring time rather than read back off the request, so the route and the
// panel it renders are bound in one place and a URL cannot name a panel that
// does not exist.
//
// The root keeps honouring the `?cd=` links the site used to make, when every
// directory was a signal on one page. Those are out in the world -- a post
// page linked back that way -- so they redirect to the route the directory
// lives at now rather than quietly landing on the listing.
func (h *Handlers) PortfolioPage(tab string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tab == components.TabHome {
			if cd := r.URL.Query().Get("cd"); cd != "" {
				for _, d := range components.Dirs() {
					if d.Name == cd || d.Key == cd {
						http.Redirect(w, r, d.Href(), http.StatusMovedPermanently)
						return
					}
				}
			}
		}

		filter := r.URL.Query().Get("filter")

		if err := pages.PortfolioPage(tab, filter).Render(r.Context(), w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// NotFoundHandler is the site's answer to a path nothing serves. It lives with
// the portfolio because the answer is the front page's listing -- what a
// visitor needs after a wrong turn is the set of directories that do exist,
// and that is the one page that has it.
func NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if err := pages.NotFound(components.NotFoundPath(r.URL.Path)).Render(r.Context(), w); err != nil {
			return
		}
	}
}
