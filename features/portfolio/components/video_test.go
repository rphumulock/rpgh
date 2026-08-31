package components

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

// iframeTag is the player element inside a rendered card.
var iframeTag = regexp.MustCompile(`(?s)<iframe.*?>`)

func renderCard(t *testing.T, v Video) string {
	t.Helper()
	var b strings.Builder
	if err := VideoCard(v).Render(context.Background(), &b); err != nil {
		t.Fatalf("rendering the card: %v", err)
	}
	return b.String()
}

// TestPlayerSendsAnOrigin guards a failure nothing else here can see. YouTube
// reads the embedding origin off the Referer header; an iframe sending none
// gets a player that renders "Video player configuration error" instead of the
// video, and the server, the tests and the CSP all still look correct. The
// still beside it is deliberately no-referrer, so the two are easy to conflate.
func TestPlayerSendsAnOrigin(t *testing.T) {
	html := renderCard(t, Video{Episode: "X", Title: "t", ID: "xq1dVQ-isb4"})
	iframe := iframeTag.FindString(html)
	if iframe == "" {
		t.Fatal("the card rendered no iframe, so there is no player to click into")
	}
	if strings.Contains(iframe, `referrerpolicy="no-referrer"`) {
		t.Error("the player sends no referrer, so YouTube will refuse to configure it")
	}
	if !strings.Contains(iframe, `referrerpolicy="strict-origin-when-cross-origin"`) {
		t.Errorf("the player does not name a referrer policy that sends an origin: %s", iframe)
	}
}

// TestPlayerHasNoSrcUntilPlayed is the other half of the facade: a src on the
// element as rendered would fetch YouTube's player for every visitor who opens
// the videos tab, which is the thing the card is shaped the way it is to avoid.
func TestPlayerHasNoSrcUntilPlayed(t *testing.T) {
	iframe := iframeTag.FindString(renderCard(t, Video{Episode: "X", Title: "t", ID: "xq1dVQ-isb4"}))
	if strings.Contains(iframe, " src=") {
		t.Errorf("the player carries a src before it is played: %s", iframe)
	}
	if !strings.Contains(iframe, "data-attr:src=") {
		t.Errorf("the player never binds a src, so clicking it does nothing: %s", iframe)
	}
}
