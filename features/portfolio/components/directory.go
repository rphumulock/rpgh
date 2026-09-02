package components

import (
	"net/url"
	"strings"

	"rpgh/features/blog/content"
)

// The panels the site is made of. A key names a panel in Go; what a visitor
// travels by is the route under Href, one per directory. The two are separate
// because the key is an identifier and the route is part of the address bar --
// `tech/` is spelled `stack` here and has been since before it was a route.
const (
	TabHome     = "home"
	TabProjects = "projects"
	TabStack    = "stack"
	TabVideos   = "videos"
	TabBlogs    = "blogs"
)

// Dir is one row of the front page listing: a directory a visitor can cd into.
// Count and Unit are what the row reports in place of a file size, since the
// interesting number about a directory here is how much is inside it.
type Dir struct {
	Key   string
	Name  string
	Icon  string
	Blurb string
	Count int
	Unit  string
}

// Mode is the permission column. Every entry here is a directory and none of
// them is writable by a visitor, so the string is the same on every row -- it
// is column dressing, and stating it once keeps it that way.
const Mode = "drwxr-xr-x"

// Size is how a row reports what it holds, standing in for the byte count an
// `ls -l` line would carry there.
func (d Dir) Size() string {
	return plural(d.Count, d.Unit)
}

// Root is the directory the listing lists, and the one `../` climbs back to.
// HomeHref is that same place as a URL: the front page is the root directory,
// so climbing out of a panel is a link to it rather than a signal change.
const (
	Root     = "~/rpgh"
	HomeHref = "/"
)

// Path is where a panel says you are, written the way you got there. A Dir
// with no name is the root itself, which is the listing rather than a row in
// it.
func (d Dir) Path() string {
	if d.Name == "" {
		return Root
	}
	return Root + "/" + d.Name
}

// NotFoundPath writes a URL that leads nowhere the way the path bar writes the
// ones that lead somewhere, so the 404 reads as the same shell answering. The
// argument is a request path, which is whatever a visitor typed -- it is only
// ever printed, and templ escapes it on the way out.
func NotFoundPath(urlPath string) string {
	trimmed := strings.TrimSuffix(urlPath, "/")
	if trimmed == "" || trimmed == "/" {
		return Root
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return Root + trimmed
}

// Href is the route the directory is served at. It is the name a visitor sees
// in the path bar, so `~/rpgh/tech` is at /tech -- one spelling, in the URL
// and on the page both.
func (d Dir) Href() string {
	if d.Name == "" {
		return HomeHref
	}
	return "/" + d.Name
}

// DirByTab finds the row that opens a panel, so a panel can name itself from
// the same entry the listing rendered rather than spelling the name twice.
func DirByTab(key string) Dir {
	for _, d := range Dirs() {
		if d.Key == key {
			return d
		}
	}
	return Dir{Key: key}
}

// Dirs is the front page. Counts are read off the same data the panels render,
// so a project added below shows up in the listing without a second edit.
func Dirs() []Dir {
	return []Dir{
		{
			Key:   TabProjects,
			Name:  "projects",
			Icon:  "lucide:folder-git-2",
			Blurb: "Things I built end to end, filterable by what they were built with.",
			Count: len(Projects),
			Unit:  "project",
		},
		{
			Key:   TabStack,
			Name:  "tech",
			Icon:  "lucide:list-tree",
			Blurb: "The toolbox as a tree -- languages, databases, backend, infra, tooling.",
			Count: StackCount(),
			Unit:  "file",
		},
		{
			Key:   TabBlogs,
			Name:  content.Name,
			Icon:  "lucide:file-pen-line",
			Blurb: "Notes on what I am building, written while the reasons are still fresh.",
			Count: len(content.Posts()),
			Unit:  "post",
		},
		{
			Key:   TabVideos,
			Name:  "videos",
			Icon:  "lucide:clapperboard",
			Blurb: "Walkthroughs I have recorded, grouped by the series they belong to.",
			Count: TotalVideoCount(),
			Unit:  "video",
		},
	}
}

// ProjectsFilterHref opens the projects directory with one filter already
// applied, which is what a tech chip does: it is a link to the same panel the
// listing links to, carrying the answer to "built with what" in the query
// rather than in a signal a page load would forget.
func ProjectsFilterHref(name string) string {
	d := DirByTab(TabProjects)
	if name == "" || name == FilterAll {
		return d.Href()
	}
	return d.Href() + "?filter=" + url.QueryEscape(name)
}
