// Package components holds the resume content. Everything here is transcribed
// from web/resources/static/assets/richard-humulock-resume.pdf -- when the PDF
// is regenerated, update this file to match so the page and the download agree.
package components

// PDFPath is the downloadable copy of this resume, relative to the static dir.
const PDFPath = "assets/richard-humulock-resume.pdf"

// Header is the block above the fold: name, title, and how to reach him.
const (
	Name      = "Richard Humulock"
	Tagline   = "Senior Software Engineer · Distributed Systems & Platform Engineering"
	Location  = "Baltimore, MD"
	Clearance = "Active Top Secret Clearance"
	Summary   = "Software engineer with 10+ years building scalable, data-driven systems end to end — " +
		"distributed backends, event-driven architecture, and cloud infrastructure. Depth in Go and JVM " +
		"services, messaging and streaming, and PostgreSQL performance, with a record of setting security " +
		"and data-governance standards in cleared environments."
)

// ContactLink is one entry in the header's contact row.
type ContactLink struct {
	Label string
	Href  string
	Icon  string
}

var Contact = []ContactLink{
	{"rphumulock@gmail.com", "mailto:rphumulock@gmail.com", "lucide:mail"},
	{"richardhumulock.fly.dev", "https://richardhumulock.fly.dev/", "lucide:globe"},
	{"linkedin.com/in/richard-humulock", "https://www.linkedin.com/in/richard-humulock/", "simple-icons:linkedin"},
	{"github.com/rphumulock", "https://github.com/rphumulock", "simple-icons:github"},
}

// Role is one position under EXPERIENCE.
type Role struct {
	Title    string
	Company  string
	Location string
	Dates    string
	Bullets  []string
}

// Roles is in reverse-chronological order, matching the PDF.
//
// NOTE: the source PDF carries literal "[START DATE]" and "[END DATE]"
// placeholders for the two most recent roles; the known halves are kept here
// and the missing halves are left out rather than guessed.
var Roles = []Role{
	{
		Title:    "Technical Marketing Engineer",
		Company:  "Synadia",
		Location: "Baltimore, MD (Remote)",
		Dates:    "Present",
		Bullets: []string{
			"Design and build reference architectures and production-grade demo applications on NATS.io and Synadia Cloud, demonstrating event-driven patterns including pub/sub, JetStream persistence, and key-value and object stores.",
			"Author developer documentation, technical guides, and tutorials for adopting NATS in production.",
			"Work across the stack in Go, TypeScript, and hypermedia frameworks to demonstrate low-latency application patterns from broker to browser.",
		},
	},
	{
		Title:    "Senior Research Engineer",
		Company:  "Two Six Technologies",
		Location: "Baltimore, MD (Remote)",
		Dates:    "September 2020",
		Bullets: []string{
			"Architected scalable backend services with Vert.x and HyperExpress, deploying a Hazelcast-backed Vert.x cluster for resilience, redundancy, and distributed state management.",
			"Implemented secure authentication workflows using OAuth2 and PKCE, and designed and enforced Role-Based Access Control for granular, auditable access to protected resources.",
			"Established data governance practices in compliance with ISM standards, defining the controls the team built against.",
			"Optimized PostgreSQL schemas and query performance; integrated WebSockets into Ember.js for real-time updates.",
			"Mentored junior engineers and raised the team baseline on security and scalability practices.",
		},
	},
	{
		Title:    "Software Development Engineer",
		Company:  "U.S. Engineering Solutions",
		Location: "Baltimore, MD",
		Dates:    "March 2018 – September 2020",
		Bullets: []string{
			"Built real-time weather monitoring dashboards in React and Redux, reducing data load times by 30%.",
			"Developed Spring Boot and Hibernate services backed by Redis to cache and serve high-frequency weather data.",
			"Designed an AWS ingestion pipeline using Lambda for real-time processing and S3 for durable, scalable storage.",
		},
	},
}

// EarlierRoles is the condensed "Earlier" line that closes the experience section.
const EarlierRoles = "Eutality Group (2017–2018) · Mannakee Group (2016–2017) · e-Management (2014–2016) — " +
	"C#/.NET physics simulation, AngularJS and SharePoint applications for high-security environments, " +
	"and a Java web crawler with D3."

// Project is one entry under SELECTED PROJECTS.
type Project struct {
	Name  string
	Tech  string
	Blurb string
	Href  string
}

var Projects = []Project{
	{
		Name:  "Distributed Sensor Cluster",
		Tech:  "Vert.x · Hazelcast · Datastar · Docker Compose",
		Blurb: "Three-node Vert.x cluster where virtual sensors are deployed and torn down at runtime. Nodes coordinate over the clustered event bus with Hazelcast managing membership; readings stream to the UI over SSE rather than polling.",
		Href:  "https://richardhumulock.fly.dev/projects/vertx_hazelcast_cluster",
	},
	{
		Name:  "Real-Time Multiplayer over Redis Pub/Sub",
		Tech:  "Node.js · WebSockets · Redis",
		Blurb: "Implemented the WebSocket handshake by hand rather than using a library, broadcasting moves through Redis Pub/Sub; game state persists so matches survive a server restart.",
		Href:  "https://github.com/rphumulock/redis_websockets_tictactoe",
	},
	{
		Name:  "ECS Game Server",
		Tech:  "Go · entity-component-system",
		Blurb: "ECS design separating data from behavior; systems iterate contiguous component arrays for cache locality.",
		Href:  "https://richardhumulock.fly.dev/projects/go_ecs_roguelike",
	},
}

// SkillGroup is one labelled row under TECHNICAL SKILLS.
type SkillGroup struct {
	Label string
	Items string
}

var Skills = []SkillGroup{
	{"Languages", "Go, TypeScript, JavaScript, Java, C#, SQL"},
	{"Messaging", "NATS.io / JetStream, RabbitMQ, Redis Streams, Hazelcast, WebSockets, gRPC"},
	{"Architecture", "Event-driven systems, microservices, event sourcing, CQRS, SAGA, serverless"},
	{"Backend", "Vert.x, Spring Boot, Node.js (Express, Fastify), ASP.NET Core, REST, GraphQL"},
	{"Data", "PostgreSQL, Redis, MongoDB, MySQL, schema design, query and index optimization"},
	{"Frontend", "React, Vue.js, Svelte, Angular, Ember.js, Datastar, TailwindCSS"},
	{"Cloud & DevOps", "AWS (Lambda, S3), Docker, Kubernetes / K3s, Fly.io, GitHub Actions, GitLab CI"},
	{"Security", "OAuth2 / PKCE, JWT, Role-Based Access Control, ISM compliance"},
	{"Testing", "Playwright, k6, Jest, JUnit"},
}

// School is one entry under EDUCATION.
type School struct {
	Degree string
	Name   string
	Dates  string
}

var Education = []School{
	{"B.S., Computer Science", "University of Maryland, College Park, MD", "2010 – 2014"},
	{"A.S., Computer Science", "College of Southern Maryland, La Plata, MD", "2008 – 2010"},
}
