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

// ShowExpr is the Datastar expression that keeps a card visible for the
// current $filter signal.
func (p Project) ShowExpr() string {
	tech, err := json.Marshal(p.Tech)
	if err != nil {
		return "true"
	}
	return "$filter === 'all' || " + string(tech) + ".includes($filter)"
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
