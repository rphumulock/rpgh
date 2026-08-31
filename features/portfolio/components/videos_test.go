package components

import (
	"regexp"
	"strings"
	"testing"
)

// youTubeID is the shape of a watch id. A typo here renders a card whose
// thumbnail 404s and whose link goes nowhere, and nothing else would catch it.
var youTubeID = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

func TestVideoIDsAreWellFormed(t *testing.T) {
	seen := map[string]string{}
	for _, c := range Channels {
		for _, s := range c.Series {
			for _, v := range s.Videos {
				if !youTubeID.MatchString(v.ID) {
					t.Errorf("%s: id %q is not a YouTube watch id", v.Label(), v.ID)
				}
				if prev, dup := seen[v.ID]; dup {
					t.Errorf("%s: id %q already used by %s", v.Label(), v.ID, prev)
				}
				seen[v.ID] = v.Label()
			}
		}
	}
}

// TestPlaylistKeysAreUnique guards the $playlist signal: two containers
// sharing a key would open together and their tabs would both light up.
func TestPlaylistKeysAreUnique(t *testing.T) {
	panels := PlaylistPanels()
	if len(panels) == 0 {
		t.Fatal("no playlists, so the videos tab renders an empty section")
	}
	seen := map[string]bool{}
	for _, p := range panels {
		if p.Key() == "" {
			t.Errorf("series %q on %s has neither a playlist id nor a video to key on", p.Series.Name, p.Handle)
		}
		if seen[p.Key()] {
			t.Errorf("panel key %q is declared twice", p.Key())
		}
		seen[p.Key()] = true
	}
	if !seen[DefaultPlaylist()] {
		t.Errorf("DefaultPlaylist() = %q, which no container matches", DefaultPlaylist())
	}
}

// TestEveryPanelHasACard guards what each container shows: an empty series
// renders no cards at all, or one with a dead link and a 404 still.
func TestEveryPanelHasACard(t *testing.T) {
	for _, p := range PlaylistPanels() {
		if len(p.Series.Cards()) == 0 {
			t.Errorf("series %q holds no videos, so its container renders nothing", p.Series.Name)
		}
		for _, v := range p.Series.Cards() {
			if v.ID == "" {
				t.Errorf("series %q cards a video with no watch id", p.Series.Name)
			}
		}
	}
}

// TestCardsFollowThePlaylist pins the split the container renders on: a series
// with a playlist holds its other episodes back behind the playlist link, and
// one without has no such link, so holding anything back would hide it.
func TestCardsFollowThePlaylist(t *testing.T) {
	withList := VideoSeries{Playlist: "PL123", Videos: []Video{
		{Episode: "01", Title: "first", ID: "aaaaaaaaaaa"},
		{Episode: "02", Title: "second", ID: "bbbbbbbbbbb"},
	}}
	if got := withList.Cards(); len(got) != 1 || got[0].Title != "first" {
		t.Errorf("Cards() = %+v, want only the featured episode", got)
	}

	loose := VideoSeries{Videos: []Video{
		{Title: "one-off", ID: "aaaaaaaaaaa"},
		{Title: "another", ID: "bbbbbbbbbbb"},
	}}
	if got := loose.Cards(); len(got) != 2 {
		t.Errorf("Cards() = %+v, want every video, since nothing links to the rest", got)
	}
}

// TestLooseSeriesKeyOffTheirVideo guards the fallback the $playlist signal
// leans on: a series with no playlist id still needs a key its tab can set and
// its container can match, or the tab would open nothing.
func TestLooseSeriesKeyOffTheirVideo(t *testing.T) {
	p := PlaylistPanel{Series: VideoSeries{Videos: []Video{{Title: "one-off", ID: "aaaaaaaaaaa"}}}}
	if p.Key() != "aaaaaaaaaaa" {
		t.Errorf("Key() = %q, want the featured video's watch id", p.Key())
	}
}

func TestEpisodesAreUniqueWithinASeries(t *testing.T) {
	for _, c := range Channels {
		if c.Handle == "" || len(c.Series) == 0 {
			t.Errorf("channel %+v renders an empty container", c)
		}
		for _, s := range c.Series {
			seen := map[string]bool{}
			for _, v := range s.Videos {
				if v.Title == "" {
					t.Errorf("series %q has a video with no title: %+v", s.Name, v)
				}
				// A one-off carries no episode number, so only numbered
				// episodes are held to being numbered once.
				if v.Episode == "" {
					continue
				}
				if seen[v.Episode] {
					t.Errorf("series %q lists episode %q twice", s.Name, v.Episode)
				}
				seen[v.Episode] = true
			}
		}
	}
}

// TestAtMostOneFeaturedPerSeries guards the flag Featured() reads: two flagged
// episodes would silently render the earlier one and quietly ignore the other.
func TestAtMostOneFeaturedPerSeries(t *testing.T) {
	for _, c := range Channels {
		for _, s := range c.Series {
			var flagged []string
			for _, v := range s.Videos {
				if v.Featured {
					flagged = append(flagged, v.Label())
				}
			}
			if len(flagged) > 1 {
				t.Errorf("series %q flags %d episodes as featured, but only one renders: %v",
					s.Name, len(flagged), flagged)
			}
		}
	}
}

// TestFeaturedHonoursTheFlag pins the fallback as well as the flag: an
// unflagged series must still feature its first episode rather than nothing.
func TestFeaturedHonoursTheFlag(t *testing.T) {
	flagged := VideoSeries{Videos: []Video{
		{Episode: "01", Title: "first", ID: "aaaaaaaaaaa"},
		{Episode: "X", Title: "flagged", ID: "bbbbbbbbbbb", Featured: true},
	}}
	if got := flagged.Featured(); got.Title != "flagged" {
		t.Errorf("Featured() = %q, want the flagged episode", got.Title)
	}

	unflagged := VideoSeries{Videos: []Video{
		{Episode: "01", Title: "first", ID: "aaaaaaaaaaa"},
		{Episode: "02", Title: "second", ID: "bbbbbbbbbbb"},
	}}
	if got := unflagged.Featured(); got.Title != "first" {
		t.Errorf("Featured() = %q, want the first episode", got.Title)
	}
}

// TestEmbedIsAllowedByThePolicy ties the embed host to the CSP that has to
// allow it: changing one without the other renders a card whose player is
// blocked, and only the browser console would say so.
func TestEmbedIsAllowedByThePolicy(t *testing.T) {
	for _, p := range PlaylistPanels() {
		for _, v := range p.Series.Cards() {
			if !strings.HasPrefix(v.EmbedURL(), "https://www.youtube-nocookie.com/embed/") {
				t.Errorf("%s: embed URL %q is not on the frame-src host", v.Label(), v.EmbedURL())
			}
			if !strings.Contains(v.EmbedURL(), v.ID) {
				t.Errorf("%s: embed URL %q does not carry the watch id", v.Label(), v.EmbedURL())
			}
		}
	}
}
