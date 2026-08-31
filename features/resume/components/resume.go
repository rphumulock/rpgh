// Package components holds the resume content. Everything here is transcribed
// from web/resources/static/assets/peter-humulock-resume.pdf -- when the PDF
// is regenerated, update this file to match so the page and the download agree.
package components

// PDFPath is the downloadable copy of this resume, relative to the static dir.
const PDFPath = "assets/peter-humulock-resume.pdf"

// Header is the block above the fold: name, title, and how to reach him.
const (
	Name      = "Peter Humulock"
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
	{"rpgh.dev", "https://rpgh.dev/", "lucide:globe"},
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
var Roles = []Role{
	{
		Title:    "Technical Marketing Engineer",
		Company:  "Synadia",
		Location: "Baltimore, MD (Remote)",
		Dates:    "April 2025 – August 2026",
		Bullets: []string{
			"Designed and built reference architectures and production-grade demo applications on NATS.io and Synadia Cloud, demonstrating event-driven patterns including pub/sub, JetStream persistence, and key-value and object stores.",
			"Produced the developer documentation, technical guides, and video demos teams used to adopt NATS in production.",
			"Worked across the stack in Go, TypeScript, and hypermedia frameworks to demonstrate low-latency application patterns from broker to browser.",
		},
	},
	{
		Title:    "Senior Research Engineer",
		Company:  "Two Six Technologies",
		Location: "Baltimore, MD (Remote)",
		Dates:    "September 2020 – April 2025",
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
		Name:  "Vert.x Hazelcast Sensor Cluster",
		Tech:  "Java 21 · Vert.x · Hazelcast · Datastar · Docker",
		Blurb: "A single clustered event bus spanning three JVMs, with a Hazelcast AsyncMap as the only data store behind cluster-wide locks. Sensors deploy and tear down at runtime; readings reach the browser over SSE.",
		Href:  "https://github.com/rphumulock/vertx-ds-hazelcast",
	},
	{
		Name:  "Datastar NATS Tic-Tac-Toe",
		Tech:  "Go · NATS JetStream KV · SSE · CQRS",
		Blurb: "The server owns all state: moves are validated into JetStream KV and rendered markup is pushed to both players over SSE — no client-side application logic, no JSON API. CQRS separates writes from KV-watch reads; presence heartbeats expire orphaned games.",
		Href:  "https://github.com/rphumulock/datastar-nats-tictactoe",
	},
	{
		// The PDF gives this one a descriptor where the others carry a tech
		// list, and no blurb under it -- Tech is the line after the name in
		// both cases, so it holds the descriptor here.
		Name: "Datastar Video Series",
		Tech: "Eleven-episode YouTube series on hypermedia application patterns",
		Href: "https://www.youtube.com/playlist?list=PLbqyjFEQew904tnpc7dtc6VuyX7HikBfR",
	},
}

// SkillGroup is one labelled row under TECHNICAL SKILLS.
type SkillGroup struct {
	Label string
	Items string
}

var Skills = []SkillGroup{
	{"Languages", "Go, Java, JavaScript, TypeScript, C#, SQL"},
	{"Messaging", "NATS.io / JetStream, RabbitMQ, Redis Streams, Hazelcast, WebSockets, gRPC"},
	{"Architecture", "Event-driven systems, microservices, event sourcing, CQRS, SAGA, serverless"},
	{"Backend", "Vert.x, Spring Boot, Node.js (Express, Fastify), ASP.NET Core, REST, GraphQL"},
	{"Data", "PostgreSQL, Redis, MongoDB, MySQL, schema design, query and index optimization"},
	{"Frontend & Mobile", "React, Vue.js, Svelte, Angular, Ember.js, Datastar, TailwindCSS, Swift / iOS, Android"},
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
