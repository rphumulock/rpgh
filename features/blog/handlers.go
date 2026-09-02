package blog

import (
	"net/http"

	"rpgh/features/blog/content"
	"rpgh/features/blog/pages"

	"github.com/go-chi/chi/v5"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// PostPage serves one post. A slug nobody wrote gets a 404 carrying the same
// chrome as everything else rather than chi's bare one -- it is a page a
// visitor can arrive at from a stale link, so it should let them back in.
func (h *Handlers) PostPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	post, ok := content.BySlug(slug)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		if err := pages.NotFound(slug).Render(r.Context(), w); err != nil {
			return
		}
		return
	}

	if err := pages.PostPage(post).Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
