package portfolio

import (
	"net/http"

	"rpgh/features/portfolio/pages"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) PortfolioPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.PortfolioPage().Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
