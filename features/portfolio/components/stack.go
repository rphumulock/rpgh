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
			{"Rust", "simple-icons:rust"},
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
			{"Node.js", "simple-icons:nodedotjs"},
			{"Express", "simple-icons:express"},
			{"Vert.x", "simple-icons:eclipsevertdotx"},
			{"Hazelcast", "lucide:network"},
			{"chi", "lucide:route"},
			{"gRPC", "simple-icons:grpc"},
			{"REST", "lucide:webhook"},
			{"MCP", "simple-icons:modelcontextprotocol"},
			{"Server-Sent Events", "lucide:radio"},
			{"WebSockets", "lucide:plug-zap"},
			{"Socket.IO", "simple-icons:socketdotio"},
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
		Name: "observability",
		Items: []StackItem{
			{"Prometheus", "simple-icons:prometheus"},
			{"Grafana", "simple-icons:grafana"},
			{"OpenTelemetry", "simple-icons:opentelemetry"},
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
			{"Cargo", "simple-icons:rust"},
			{"tree-sitter", "lucide:list-tree"},
			{"Vite", "simple-icons:vite"},
		},
	},
}

// SetupItem is one line in the "what I work in" panel beside the tree. That
// panel is about the machine this all gets written on, which is a different
// question from the Stack above -- nothing here is filterable, because none of
// it is what a project was built with.
type SetupItem struct {
	Role string
	Name string
	Icon string
	Note string
}

// SetupIntro sits above the list.
const SetupIntro = "A fresh box gets here from one bash script -- same setup on every machine, " +
	"rebuilt from scratch rather than carried around."

var Setup = []SetupItem{
	{"os", "Pop!_OS", "simple-icons:popos", "Ubuntu underneath, tiling on top."},
	{"editor", "LazyVim", "simple-icons:lazyvim", "Neovim with LazyVim over it; where most of the day goes."},
	{"editor", "VS Code", "simple-icons:visualstudiocode", "For stepping through a debugger, and diffs I'd rather click through."},
	{"terminal", "WezTerm", "simple-icons:wezterm", "GPU-rendered, configured in Lua, same config everywhere."},
	{"shell", "Oh My Zsh", "devicon-plain:ohmyzsh", "zsh with completions, syntax highlighting and a git-aware prompt."},
	{"agents", "Herdr", "lucide:bot", "Keeps coding agents in sessions that outlive a disconnect."},
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
