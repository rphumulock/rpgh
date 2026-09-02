// Package content is the writing on this site. A post is a Markdown file under
// posts/, embedded into the binary like everything else here, so publishing is
// a commit rather than a deploy step and there is nothing to back up.
//
// Every post opens with a frontmatter block:
//
//	---
//	title: What the post is called
//	published: 2026-09-02
//	blurb: One line, shown on the listing.
//	draft: true
//	---
//
// The keys are the four above and nothing else -- this is deliberately not a
// YAML parser, so a key it does not know is an error rather than a field that
// silently does nothing. Values are read to the end of the line; surrounding
// quotes are stripped if present.
//
// The slug is the filename, so the file is the URL: posts/why-sse.md is served
// at /blog/why-sse.
package content

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
)

// all: is what makes this compile with no posts written yet -- the pattern has
// to match something, and without the prefix an otherwise empty directory
// holding only .gitkeep matches nothing and fails the build.
//
//go:embed all:posts
var files embed.FS

// Name is the directory the front page lists this feature under, Path is that
// directory written out -- what a post page shows above its title -- and Href
// is the page it is served at, which is where a post sends its reader back to.
// The front page builds all three from its own root; TestPathMatchesTheListing
// in features/portfolio/pages is what keeps the spellings equal.
const (
	Name = "blogs"
	Path = "~/rpgh/" + Name
	Href = "/" + Name
)

// DateLayout is how a published date is written in frontmatter and rendered on
// the page. It sorts lexically as well as chronologically, which is the reason
// to insist on it rather than accept whatever time.Parse would take.
const DateLayout = "2006-01-02"

// Post is one piece of writing. HTML is the rendered body, already trusted:
// the Markdown it came from is in this repository, so it is our own output
// rather than anything a visitor supplied.
type Post struct {
	Slug      string
	Title     string
	Blurb     string
	Published time.Time
	Draft     bool
	HTML      templ.Component
}

// Date is the published date as it is written and displayed.
func (p Post) Date() string {
	return p.Published.Format(DateLayout)
}

// Href is where the post is served.
func (p Post) Href() string {
	return "/blog/" + p.Slug
}

// parsed holds the result of reading posts/ once. Parsing at startup rather
// than per request means a malformed post is a boot failure with a filename in
// it, not a page that renders wrong for a visitor and right in a test.
var parsed = sync.OnceValues(load)

// Posts is every published post, newest first. Drafts are parsed -- so a typo
// in one still fails the build -- but never returned.
func Posts() []Post {
	posts, _ := parsed()
	return posts
}

// Err reports whether posts/ could be read. SetupRoutes surfaces it, which is
// what turns a bad post into a server that refuses to start.
func Err() error {
	_, err := parsed()
	return err
}

// BySlug finds a post by its URL segment. Drafts are not reachable, so an
// unfinished post cannot be read by guessing its filename.
func BySlug(slug string) (Post, bool) {
	for _, p := range Posts() {
		if p.Slug == slug {
			return p, true
		}
	}
	return Post{}, false
}

func load() ([]Post, error) {
	entries, err := fs.ReadDir(files, "posts")
	if err != nil {
		return nil, fmt.Errorf("reading posts: %w", err)
	}

	md := goldmark.New()
	posts := make([]Post, 0, len(entries))
	seen := map[string]string{}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		raw, err := files.ReadFile(path.Join("posts", name))
		if err != nil {
			return nil, fmt.Errorf("reading posts/%s: %w", name, err)
		}

		post, body, err := parseFront(raw)
		if err != nil {
			return nil, fmt.Errorf("posts/%s: %w", name, err)
		}

		post.Slug = strings.TrimSuffix(name, ".md")
		if prev, dup := seen[post.Slug]; dup {
			return nil, fmt.Errorf("posts/%s and posts/%s share the slug %q", prev, name, post.Slug)
		}
		seen[post.Slug] = name

		var out bytes.Buffer
		if err := md.Convert(body, &out); err != nil {
			return nil, fmt.Errorf("posts/%s: rendering markdown: %w", name, err)
		}
		post.HTML = templ.Raw(out.String())

		if post.Draft {
			continue
		}
		posts = append(posts, post)
	}

	// Newest first, and by slug when two posts share a day, so the order is
	// the same on every boot rather than whatever the filesystem returned.
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Published.Equal(posts[j].Published) {
			return posts[i].Slug < posts[j].Slug
		}
		return posts[i].Published.After(posts[j].Published)
	})

	return posts, nil
}

// fence is the line that opens and closes a frontmatter block.
const fence = "---"

// parseFront splits a post into its frontmatter and its Markdown body. It is
// strict on purpose: a post is written by one person and read by everyone, so
// a missing title should stop a deploy rather than render as an empty heading.
func parseFront(raw []byte) (Post, []byte, error) {
	text := strings.ReplaceAll(string(raw), "\r\n", "\n")

	rest, ok := strings.CutPrefix(text, fence+"\n")
	if !ok {
		return Post{}, nil, fmt.Errorf("no frontmatter: the file must open with a %q line", fence)
	}
	front, body, ok := strings.Cut(rest, "\n"+fence+"\n")
	if !ok {
		return Post{}, nil, fmt.Errorf("frontmatter is never closed by a %q line", fence)
	}

	var post Post
	for i, line := range strings.Split(front, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return Post{}, nil, fmt.Errorf("frontmatter line %d is not `key: value`: %q", i+1, line)
		}
		key = strings.TrimSpace(key)
		value = unquote(strings.TrimSpace(value))

		switch key {
		case "title":
			post.Title = value
		case "blurb":
			post.Blurb = value
		case "published":
			when, err := time.Parse(DateLayout, value)
			if err != nil {
				return Post{}, nil, fmt.Errorf("published %q is not a %s date", value, DateLayout)
			}
			post.Published = when
		case "draft":
			post.Draft = value == "true"
		default:
			return Post{}, nil, fmt.Errorf("unknown frontmatter key %q", key)
		}
	}

	switch {
	case post.Title == "":
		return Post{}, nil, fmt.Errorf("no title")
	case post.Blurb == "":
		return Post{}, nil, fmt.Errorf("no blurb, which is what the listing shows")
	case post.Published.IsZero():
		return Post{}, nil, fmt.Errorf("no published date")
	}

	return post, []byte(strings.TrimPrefix(body, "\n")), nil
}

// unquote strips one matching pair of surrounding quotes, so a title with a
// leading # or a trailing colon can be written without the parser guessing.
func unquote(s string) string {
	for _, q := range []string{`"`, `'`} {
		if len(s) >= 2 && strings.HasPrefix(s, q) && strings.HasSuffix(s, q) {
			return s[1 : len(s)-1]
		}
	}
	return s
}
