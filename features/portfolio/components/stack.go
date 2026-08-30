package components

// StackItem is one tool, language, or platform. Icon is an Iconify name
// (https://icon-sets.iconify.design); leave it empty for a bare label.
type StackItem struct {
	Name string
	Icon string
}

// StackCategory is one "directory" in the tree rendering of the stack.
type StackCategory struct {
	Name  string
	Items []StackItem
}

// Stack is the source of truth for the "$ tree ./stack" section.
var Stack = []StackCategory{
	{
		Name: "languages",
		Items: []StackItem{
			{"Go", "simple-icons:go"},
			{"Java", "simple-icons:openjdk"},
			{"TypeScript", "simple-icons:typescript"},
			{"JavaScript", "simple-icons:javascript"},
			{"Python", "simple-icons:python"},
			{"C#", "simple-icons:dotnet"},
			{"SQL", "lucide:database"},
			{"Bash", "simple-icons:gnubash"},
			{"HTML", "simple-icons:html5"},
			{"CSS", "simple-icons:css"},
		},
	},
	{
		Name: "databases",
		Items: []StackItem{
			{"PostgreSQL", "simple-icons:postgresql"},
			{"SQLite", "simple-icons:sqlite"},
			{"NATS JetStream", "simple-icons:natsdotio"},
			{"Redis", "simple-icons:redis"},
			{"MongoDB", "simple-icons:mongodb"},
			{"DuckDB", "simple-icons:duckdb"},
			{"EF Core", "simple-icons:dotnet"},
		},
	},
	{
		Name: "backend",
		Items: []StackItem{
			{"NATS", "simple-icons:natsdotio"},
			{"Vert.x", "simple-icons:eclipsevertdotx"},
			{"Hazelcast", "lucide:network"},
			{"chi", "lucide:route"},
			{"gRPC", "simple-icons:grpc"},
			{"REST", "lucide:webhook"},
			{"Server-Sent Events", "lucide:radio"},
			{"WebSockets", "lucide:plug-zap"},
			{"OpenTelemetry", "simple-icons:opentelemetry"},
		},
	},
	{
		Name: "frontend",
		Items: []StackItem{
			{"Datastar", "lucide:zap"},
			{"templ", "lucide:file-code-2"},
			{"Tailwind CSS", "simple-icons:tailwindcss"},
			{"daisyUI", "simple-icons:daisyui"},
			{"Lit", "simple-icons:lit"},
			{"htmx", "simple-icons:htmx"},
			{"Alpine.js", "simple-icons:alpinedotjs"},
			{"React", "simple-icons:react"},
		},
	},
	{
		Name: "infra",
		Items: []StackItem{
			{"Docker", "simple-icons:docker"},
			{"Linux", "simple-icons:linux"},
			{"Raspberry Pi", "simple-icons:raspberrypi"},
			{"GitHub Actions", "simple-icons:githubactions"},
			{"Caddy", "simple-icons:caddy"},
			{"nginx", "simple-icons:nginx"},
			{"Cloudflare", "simple-icons:cloudflare"},
			{"Terraform", "simple-icons:terraform"},
		},
	},
	{
		Name: "tooling",
		Items: []StackItem{
			{"Git", "simple-icons:git"},
			{"Neovim", "simple-icons:neovim"},
			{"VS Code", "simple-icons:visualstudiocode"},
			{"Task", "lucide:list-checks"},
			{"Air", "lucide:refresh-cw"},
			{"esbuild", "simple-icons:esbuild"},
			{"Delve", "lucide:bug"},
			{"Maven", "simple-icons:apachemaven"},
			{"Selenium", "simple-icons:selenium"},
		},
	},
}

// TreeGlyph returns the box-drawing connector for a category at index i,
// so the section reads like real `tree` output.
func TreeGlyph(i int) string {
	if i == len(Stack)-1 {
		return "└──"
	}
	return "├──"
}

// TreePipe returns the continuation column drawn beneath a category: a pipe
// while more categories follow, blank once the tree has closed.
func TreePipe(i int) string {
	if i == len(Stack)-1 {
		return "   "
	}
	return "│  "
}

// StackCount is the total number of items, used in the section footer.
func StackCount() int {
	n := 0
	for _, c := range Stack {
		n += len(c.Items)
	}
	return n
}
