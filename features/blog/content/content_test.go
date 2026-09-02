package content

import (
	"strings"
	"testing"
	"time"
)

const good = `---
title: Why SSE
published: 2026-09-02
blurb: One line about it.
---

Body text.
`

func TestParsesAPost(t *testing.T) {
	post, body, err := parseFront([]byte(good))
	if err != nil {
		t.Fatalf("parsing a well-formed post: %v", err)
	}
	if post.Title != "Why SSE" {
		t.Errorf("title = %q", post.Title)
	}
	if post.Blurb != "One line about it." {
		t.Errorf("blurb = %q", post.Blurb)
	}
	if want := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC); !post.Published.Equal(want) {
		t.Errorf("published = %v, want %v", post.Published, want)
	}
	if got := strings.TrimSpace(string(body)); got != "Body text." {
		t.Errorf("body = %q, want the Markdown with the frontmatter stripped", got)
	}
}

// TestRejectsBadFrontmatter is the whole reason parsing happens at boot: each
// of these renders as a plausible-looking page if it is allowed through, and
// the one who notices is a reader rather than the author.
func TestRejectsBadFrontmatter(t *testing.T) {
	cases := map[string]string{
		"no frontmatter":    "Body text.\n",
		"never closed":      "---\ntitle: T\n\nBody text.\n",
		"no title":          "---\npublished: 2026-09-02\nblurb: b\n---\n\nBody.\n",
		"no blurb":          "---\ntitle: T\npublished: 2026-09-02\n---\n\nBody.\n",
		"no date":           "---\ntitle: T\nblurb: b\n---\n\nBody.\n",
		"date wrong layout": "---\ntitle: T\nblurb: b\npublished: Sept 2 2026\n---\n\nBody.\n",
		"unknown key":       "---\ntitle: T\nblurb: b\npublished: 2026-09-02\nauthor: me\n---\n\nBody.\n",
		"not a pair":        "---\ntitle: T\nblurb: b\npublished: 2026-09-02\njust a line\n---\n\nBody.\n",
	}

	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			if _, _, err := parseFront([]byte(raw)); err == nil {
				t.Error("parsed without complaint, so this would ship")
			}
		})
	}
}

// TestTitleMayCarryAColon guards the one place the line-splitting is naive:
// cutting on the first colon is right, and a quoted value has to survive it.
func TestTitleMayCarryAColon(t *testing.T) {
	raw := "---\ntitle: \"Datastar: a tour\"\nblurb: b\npublished: 2026-09-02\n---\n\nBody.\n"
	post, _, err := parseFront([]byte(raw))
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if post.Title != "Datastar: a tour" {
		t.Errorf("title = %q", post.Title)
	}
}

// TestDraftsAreNotPublished covers both halves of what draft means: kept out of
// the listing, and unreachable by guessing the URL.
func TestDraftsAreNotPublished(t *testing.T) {
	post, _, err := parseFront([]byte("---\ntitle: T\nblurb: b\npublished: 2026-09-02\ndraft: true\n---\n\nBody.\n"))
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if !post.Draft {
		t.Fatal("draft: true did not mark the post as a draft")
	}
	for _, p := range Posts() {
		if p.Draft {
			t.Errorf("draft %q is in the published listing", p.Slug)
		}
		if _, ok := BySlug(p.Slug); !ok && p.Draft {
			t.Errorf("draft %q is reachable by slug", p.Slug)
		}
	}
}

// TestPostsDirectoryParses is what makes a typo in a real post a failed build
// rather than a failed boot on the server.
func TestPostsDirectoryParses(t *testing.T) {
	if err := Err(); err != nil {
		t.Fatalf("reading posts/: %v", err)
	}
}

// TestPostsAreNewestFirst holds even with nothing written, so the order is
// settled before there is anything to get it wrong.
func TestPostsAreNewestFirst(t *testing.T) {
	posts := Posts()
	for i := 1; i < len(posts); i++ {
		if posts[i-1].Published.Before(posts[i].Published) {
			t.Errorf("%q (%s) is listed above %q (%s)",
				posts[i-1].Slug, posts[i-1].Date(), posts[i].Slug, posts[i].Date())
		}
	}
}
