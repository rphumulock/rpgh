package components

import (
	"regexp"
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
			t.Errorf("series %q on %s has no playlist id", p.Series.Name, p.Handle)
		}
		if seen[p.Key()] {
			t.Errorf("playlist id %q is declared twice", p.Key())
		}
		seen[p.Key()] = true
	}
	if !seen[DefaultPlaylist()] {
		t.Errorf("DefaultPlaylist() = %q, which no container matches", DefaultPlaylist())
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
				if v.Episode == "" || v.Title == "" {
					t.Errorf("series %q has an episode missing a number or title: %+v", s.Name, v)
				}
				if seen[v.Episode] {
					t.Errorf("series %q lists episode %q twice", s.Name, v.Episode)
				}
				seen[v.Episode] = true
			}
		}
	}
}
