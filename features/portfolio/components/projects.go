package components

import (
	"encoding/json"
	"strings"
)

// Project is one entry in the projects grid. Image is a path relative to the
// embedded static dir, resolved through resources.StaticPath at render time.
// Tech names must match a StackItem.Name exactly -- that is what wires the
// stack tab's filter to this grid. TestTechNamesResolve guards the link.
type Project struct {
	Title string
	Blurb string
	Tech  []string
	Image string
	Repo  string
	Demo  string
}

// Projects is the source of truth for the grid. Order here is display order.
var Projects = []Project{
	{
		Title: "xsym",
		Blurb: "A cross-language structural code index, served over MCP. tree-sitter parses Go, Python and Rust into one SQLite index, so the same type spelled three ways collapses to a single lookup -- something grep cannot do, because the names genuinely differ.",
		Tech: []string{
			"Rust",
			"SQLite",
			"MCP", "tree-sitter",
			"Git", "Cargo",
		},
		Image: "assets/projects/xsym.svg",
		Repo:  "https://github.com/rphumulock/xsym",
	},
	{
		Title: "Datastar NATS Tic Tac Toe",
		Blurb: "Real-time multiplayer tic tac toe. Game state lives in NATS JetStream; moves stream to every player over SSE with no client-side framework.",
		Tech: []string{
			"Go", "HTML", "CSS",
			"NATS", "NATS JetStream", "chi", "Server-Sent Events",
			"Datastar", "templ", "Tailwind CSS", "daisyUI",
			"Docker",
			"Git", "Task", "Air", "Delve",
		},
		Image: "assets/projects/datastar-nats-tictactoe.svg",
		Repo:  "https://github.com/rphumulock/datastar-nats-tictactoe",
	},
	{
		Title: "Redis WebSockets Tic Tac Toe",
		Blurb: "The same game a stack earlier: a React front end over raw WebSockets or Socket.IO, picked by one env var, with Redis pub/sub fanning every move out to both tabs and holding the board.",
		Tech: []string{
			"JavaScript", "HTML", "CSS",
			"Node.js", "Express", "WebSockets", "Socket.IO",
			"Redis",
			"React", "Tailwind CSS",
			"Docker",
			"Git", "Vite",
		},
		Image: "assets/projects/redis-websockets-tictactoe.svg",
		Repo:  "https://github.com/rphumulock/redis_websockets_tictactoe",
	},
	{
		Title: "Vert.x Hazelcast Sensor Cluster",
		Blurb: "Three JVM nodes sharing one clustered event bus. Sensors are deployed and torn down from a browser tab onto whichever node you pick; Hazelcast holds the only state, and readings stream to every open tab over SSE.",
		Tech: []string{
			"Java", "HTML", "CSS",
			"Vert.x", "Hazelcast", "Server-Sent Events",
			"Datastar", "Tailwind CSS", "daisyUI",
			"Docker",
			"Git", "Maven",
		},
		Image: "assets/projects/vertx-ds-hazelcast.svg",
		Repo:  "https://github.com/rphumulock/vertx-ds-hazelcast",
	},
	{
		Title: "Claims Portal Crawler",
		Blurb: "A .NET crawler that files a month of claims into a JavaScript-heavy payer portal. A Postgres table is the checkpoint, not the filing cabinet: a run that dies to a stale element or an expired session resumes where it stopped.",
		Tech: []string{
			"C#", "SQL",
			"Selenium", "EF Core", "PostgreSQL",
			"Git",
		},
		Image: "assets/projects/dotnet-selenium-crawler.svg",
		Repo:  "https://github.com/rphumulock/dotnet-selenium-crawler",
	},
}

// FilterAll is the filter that hides nothing, and the value the page opens on
// unless a link asked for something narrower.
const FilterAll = "all"

// FilterOrAll answers what a page should open filtered by. A filter arrives in
// a URL, so it is whatever someone typed there: a name no project was built
// with would hide every card and read as an empty directory, and showing them
// all is the truer answer to a name that means nothing here.
func FilterOrAll(name string) string {
	if name == "" || name == FilterAll || TechUseCount(name) == 0 {
		return FilterAll
	}
	return name
}

// ShowExpr is the Datastar expression that keeps a card visible for the
// current $filter signal.
func (p Project) ShowExpr() string {
	tech, err := json.Marshal(p.Tech)
	if err != nil {
		return "true"
	}
	return "$filter === '" + FilterAll + "' || " + string(tech) + ".includes($filter)"
}

// TechUseCount is how many projects were built with a given tech. The stack
// tab uses it to decide which items are worth making clickable.
func TechUseCount(name string) int {
	n := 0
	for _, p := range Projects {
		for _, t := range p.Tech {
			if t == name {
				n++
				break
			}
		}
	}
	return n
}

// jsQuote escapes a name for embedding inside a single-quoted JavaScript
// string literal in a Datastar expression.
func jsQuote(s string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s)
}

// SetFilterExpr is the Datastar expression that selects a filter value.
func SetFilterExpr(name string) string {
	return "$filter = '" + jsQuote(name) + "'"
}

// FilterActiveExpr is the Datastar class expression marking a filter button
// active when it matches the current $filter signal.
func FilterActiveExpr(name, class string) string {
	return "{'" + class + "': $filter === '" + jsQuote(name) + "'}"
}
