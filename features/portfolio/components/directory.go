package components

import "rpgh/features/blog/content"

// The panels the site is made of. A key is the value carried in the $tab
// signal, and is what both the front page listing and the tab bar set -- they
// are two ways into the same panel, not two places.
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
const Root = "~/rpgh"

// Path is where a panel says you are, written the way you got there. A Dir
// with no name is the root itself, which is the listing rather than a row in
// it.
func (d Dir) Path() string {
	if d.Name == "" {
		return Root
	}
	return Root + "/" + d.Name
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

// SetTabExpr is the Datastar expression that opens a panel.
func SetTabExpr(key string) string {
	return "$tab = '" + jsQuote(key) + "'"
}

// TabSelectedExpr is true while a panel is the one on show.
func TabSelectedExpr(key string) string {
	return "$tab === '" + jsQuote(key) + "'"
}

// TabActiveExpr is the class expression marking the open tab, matching how the
// playlist tabs inside the videos panel are styled.
func TabActiveExpr(key string) string {
	sel := TabSelectedExpr(key)
	return "{'border-primary': " + sel + ", 'text-primary': " + sel + ", 'opacity-50': !(" + sel + ")}"
}
