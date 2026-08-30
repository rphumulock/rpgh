package resume

import (
	"net/http"

	"rpgh/features/resume/pages"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) ResumePage(w http.ResponseWriter, r *http.Request) {
	if err := pages.ResumePage().Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
