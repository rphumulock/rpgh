package host

import (
	"net/http"
	"time"

	"github.com/starfederation/datastar-go/datastar"
)

const (
	// tick is how often a stream sends. The numbers move slowly and nobody is
	// watching a footer for sub-second changes; the collector's cache means
	// the cost of a faster tick would be serialisation rather than /proc, but
	// there is no reason to pay either.
	tick = time.Second

	// maxStreams caps concurrent streams. A stream is one request that then
	// lives indefinitely, so the rate limiter never sees it again -- without a
	// cap, holding connections open is a free way to accumulate goroutines.
	maxStreams = 64

	// writeGrace bounds a single write, not the stream. Clearing the deadline
	// outright would let one stalled reader pin a goroutine forever, which is
	// exactly what the server's WriteTimeout exists to prevent; renewing it
	// per write keeps that protection while letting the stream live on.
	writeGrace = 10 * time.Second
)

type Handlers struct {
	slots chan struct{}
}

func NewHandlers() *Handlers {
	return &Handlers{slots: make(chan struct{}, maxStreams)}
}

// Events streams the host stats until the client goes away.
func (h *Handlers) Events(w http.ResponseWriter, r *http.Request) {
	select {
	case h.slots <- struct{}{}:
		defer func() { <-h.slots }()
	default:
		// Retry-After is what stops a refused client reconnecting immediately
		// and spending the capacity it is waiting for.
		w.Header().Set("Retry-After", "30")
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	rc := http.NewResponseController(w)
	sse := datastar.NewSSE(w, r)

	t := time.NewTicker(tick)
	defer t.Stop()

	for {
		if err := rc.SetWriteDeadline(time.Now().Add(writeGrace)); err != nil {
			return
		}
		if err := sse.MarshalAndPatchSignals(Collect()); err != nil {
			// The client is gone or the write timed out. Either way the only
			// thing left to do is release the slot, which the defer handles.
			return
		}

		select {
		case <-r.Context().Done():
			return
		case <-t.C:
		}
	}
}
