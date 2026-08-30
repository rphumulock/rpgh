package components

import "fmt"

// Video is one episode on the YouTube channel. ID is the watch id, which is
// all that is needed to build both the thumbnail and the watch link -- storing
// either of those verbatim would just be the same id spelled longer.
type Video struct {
	Episode   string
	Title     string
	Blurb     string
	ID        string
	Published string
}

// VideoSeries is a run of episodes that belong together, mirroring how Stack
// groups items into categories. Playlist is the YouTube playlist id.
type VideoSeries struct {
	Name     string
	Blurb    string
	Playlist string
	Videos   []Video
}

// Channel is one YouTube channel. It is the top level of the videos tab: each
// one renders as a single container holding whatever series live on it, so a
// second channel is a data edit here rather than a template change.
type Channel struct {
	Handle string
	Blurb  string
	Series []VideoSeries
}

// Channels is the source of truth for the videos tab. Episodes within a series
// are in viewing order, not upload order, which is why Published is carried
// along rather than inferred from the position.
var Channels = []Channel{
	{
		Handle: "@lockincode",
		Blurb:  "Long-form walkthroughs of the hypermedia stack I actually build on.",
		Series: []VideoSeries{
			{
				Name:     "datastar series",
				Blurb:    "A walk through Datastar from first principles: what hypermedia is, how the wire protocol works, and what every data-* attribute actually does.",
				Playlist: "PLbqyjFEQew904tnpc7dtc6VuyX7HikBfR",
				Videos: []Video{
					{
						Episode:   "01",
						Title:     "Datastar",
						Blurb:     "Series overview -- what Datastar is and where the rest of these episodes are going.",
						ID:        "I8QLWWPGT-c",
						Published: "2025-11-04",
					},
					{
						Episode:   "02",
						Title:     "Rockets Eye Overview",
						Blurb:     "The whole framework from altitude: reactive signals on the client, hypermedia on the wire.",
						ID:        "zQAz7fV95OU",
						Published: "2025-11-04",
					},
					{
						Episode:   "03",
						Title:     "Datastar Example Flow",
						Blurb:     "One request traced end to end, from the attribute that fires it to the fragment that lands.",
						ID:        "jjr9grYV8rg",
						Published: "2025-11-04",
					},
					{
						Episode:   "04",
						Title:     "Hypermedia History Lesson, Part 1",
						Blurb:     "A quick history of hypermedia and what the word actually means.",
						ID:        "jOIlaDvPhgg",
						Published: "2025-11-04",
					},
					{
						Episode:   "05",
						Title:     "Hypermedia History Lesson, Part 2",
						Blurb:     "How the web got from documents to single-page apps, and back around again.",
						ID:        "UAKJsyBhRLo",
						Published: "2025-11-04",
					},
					{
						Episode:   "06",
						Title:     "Data Attributes Overview",
						Blurb:     "The data-* attributes that make up Datastar's entire API surface.",
						ID:        "Cw3JIIHwWbM",
						Published: "2025-11-10",
					},
					{
						Episode:   "07",
						Title:     "Data Attributes Order",
						Blurb:     "When the order you write attributes in changes what they do.",
						ID:        "zOeKOf2Rk4g",
						Published: "2025-11-10",
					},
					{
						Episode:   "08",
						Title:     "Data Attributes Casing",
						Blurb:     "Where casing matters, and why the DOM is the reason it does.",
						ID:        "bZG8xknnxek",
						Published: "2025-11-10",
					},
					{
						Episode:   "09",
						Title:     "Data Attributes Aliasing",
						Blurb:     "Aliasing the data- prefix, and when that is worth doing.",
						ID:        "mX-eK6jVNIY",
						Published: "2025-11-10",
					},
					{
						Episode:   "10",
						Title:     "Data Attributes Error Handling",
						Blurb:     "The runtime error handling built into the attributes, and how to read what it tells you.",
						ID:        "JC2-D9sXR0w",
						Published: "2025-11-10",
					},
					{
						Episode:   "X",
						Title:     "Server Sent Events Overview",
						Blurb:     "A side trip through HTTP and SSE theory -- why Datastar rides on it, and why it is simpler than it sounds.",
						ID:        "xq1dVQ-isb4",
						Published: "2025-11-17",
					},
				},
			},
		},
	},
}

// WatchURL is the canonical link to an episode.
func (v Video) WatchURL() string {
	return "https://www.youtube.com/watch?v=" + v.ID
}

// ThumbURL is YouTube's 320x180 still for an episode. mqdefault is the largest
// size guaranteed to exist for every video and to be a clean 16:9 -- hqdefault
// and sddefault are 4:3 and arrive letterboxed, and hq720 is often a 404.
func (v Video) ThumbURL() string {
	return "https://i.ytimg.com/vi/" + v.ID + "/mqdefault.jpg"
}

// Label is how an episode is announced on its card.
func (v Video) Label() string {
	return fmt.Sprintf("Episode %s - %s", v.Episode, v.Title)
}

// PlaylistURL is the link to watch a series start to finish on YouTube.
func (s VideoSeries) PlaylistURL() string {
	return "https://www.youtube.com/playlist?list=" + s.Playlist
}

// URL is the channel's page. The handle is the canonical form of it now, so
// there is nothing else to store.
func (c Channel) URL() string {
	return "https://www.youtube.com/" + c.Handle
}

// PlaylistPanel is what the videos tab actually renders: one playlist, plus
// the channel it came from. Channels stays the source of truth; this is the
// flattened view of it, since a playlist is both a tab and a container there.
type PlaylistPanel struct {
	Handle     string
	ChannelURL string
	Series     VideoSeries
}

// Key is the value carried in the $playlist signal. YouTube playlist ids are
// globally unique, so two channels can never collide on one.
func (p PlaylistPanel) Key() string {
	return p.Series.Playlist
}

// PlaylistPanels flattens Channels in declaration order: every playlist of the
// first channel, then the next channel's, and so on.
func PlaylistPanels() []PlaylistPanel {
	panels := make([]PlaylistPanel, 0, len(Channels))
	for _, c := range Channels {
		for _, s := range c.Series {
			panels = append(panels, PlaylistPanel{
				Handle:     c.Handle,
				ChannelURL: c.URL(),
				Series:     s,
			})
		}
	}
	return panels
}

// DefaultPlaylist is the tab the videos section opens on. It is empty only if
// no channel declares a playlist, in which case the section renders nothing
// and the signal has nothing to match anyway.
func DefaultPlaylist() string {
	panels := PlaylistPanels()
	if len(panels) == 0 {
		return ""
	}
	return panels[0].Key()
}

// plural counts a thing for display. The footer reads "1 channel" today, so
// the s cannot just be hardcoded on.
func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// TotalVideoCount is the episode count across every channel, used in the
// section footer the way StackCount is.
func TotalVideoCount() int {
	n := 0
	for _, c := range Channels {
		for _, s := range c.Series {
			n += len(s.Videos)
		}
	}
	return n
}

// SetPlaylistExpr is the Datastar expression that selects a playlist tab.
func SetPlaylistExpr(key string) string {
	return "$playlist = '" + jsQuote(key) + "'"
}

// PlaylistSelectedExpr is true while a playlist is the one on show.
func PlaylistSelectedExpr(key string) string {
	return "$playlist === '" + jsQuote(key) + "'"
}

// PlaylistActiveExpr is the class expression marking the selected tab, matching
// how the section tabs above it are styled.
func PlaylistActiveExpr(key string) string {
	sel := PlaylistSelectedExpr(key)
	return "{'border-primary': " + sel + ", 'text-primary': " + sel + ", 'opacity-50': !(" + sel + ")}"
}
