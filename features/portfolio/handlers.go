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

// PortfolioPage renders the front page, opened on whichever directory `?cd=`
// named. The page is one document with every panel in it, so this only picks
// the starting signal -- there is no route per directory, and a name nothing
// lists falls back to the listing rather than to a blank page.
func (h *Handlers) PortfolioPage(w http.ResponseWriter, r *http.Request) {
	tab := components.TabHome
	if cd := r.URL.Query().Get("cd"); cd != "" {
		for _, d := range components.Dirs() {
			if d.Name == cd || d.Key == cd {
				tab = d.Key
				break
			}
		}
	}

	if err := pages.PortfolioPage(tab).Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
