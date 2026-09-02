package pages

import (
	"context"
	"html"
	"strings"
	"testing"

	"rpgh/features/blog/content"
	portfolio "rpgh/features/portfolio/components"
)

// render is one directory of the site as a visitor receives it. A page is the
// panel its route named and nothing else now, so which page is under test has
// to be said out loud -- that is the point of the change these tests guard.
func render(t *testing.T, tab string) string {
	t.Helper()
	var b strings.Builder
	if err := PortfolioPage(tab, "").Render(context.Background(), &b); err != nil {
		t.Fatalf("rendering %q: %v", tab, err)
	}
	return html.UnescapeString(b.String())
}

// TestEveryDirectoryRowLinksToItsPage is the guard that keeps the listing and
// the router agreed. A row is a link now, so a row pointing somewhere nothing
// is served is a 404 a visitor reaches by clicking the front page.
func TestEveryDirectoryRowLinksToItsPage(t *testing.T) {
	home := render(t, portfolio.TabHome)

	for _, d := range portfolio.Dirs() {
		if !strings.Contains(home, `href="`+d.Href()+`"`) {
			t.Errorf("the listing has no link to %q at %s", d.Name, d.Href())
		}
		if page := render(t, d.Key); !strings.Contains(page, d.Path()) {
			t.Errorf("%s renders a page that does not name %q", d.Href(), d.Path())
		}
	}
}

// TestEveryPanelCanGetBack covers the ../ in the path bar. The back button is
// the other way out and needs nothing from us; this is the one on the page,
// and a panel without it is a directory a visitor can only leave by reloading.
func TestEveryPanelCanGetBack(t *testing.T) {
	for _, d := range portfolio.Dirs() {
		if !strings.Contains(render(t, d.Key), `href="`+portfolio.HomeHref+`"`) {
			t.Errorf("%s has no ../ back to %s", d.Href(), portfolio.Root)
		}
	}
}

// TestOnlyTheNamedPanelRenders is what makes the back button worth having: a
// page is one directory. If every panel still shipped in every document, the
// routes would be four addresses for the same page.
func TestOnlyTheNamedPanelRenders(t *testing.T) {
	for _, d := range portfolio.Dirs() {
		page := render(t, d.Key)
		for _, other := range portfolio.Dirs() {
			if other.Key == d.Key {
				continue
			}
			if strings.Contains(page, other.Path()) {
				t.Errorf("%s also renders the %q panel", d.Href(), other.Name)
			}
		}
	}
}

// TestTheRootIsTheListing guards the front page being the front page: the root
// renders the listing, and no panel's path bar with it.
func TestTheRootIsTheListing(t *testing.T) {
	home := render(t, portfolio.TabHome)

	if !strings.Contains(home, "ls -l "+portfolio.Root) {
		t.Errorf("the root does not render the directory listing")
	}
	for _, d := range portfolio.Dirs() {
		if strings.Contains(home, d.Path()) {
			t.Errorf("the root also renders the %q panel", d.Name)
		}
	}
}

// TestPathMatchesTheListing keeps the two spellings of ~/rpgh/blogs equal, and
// the two spellings of the route it is served at. The blog cannot import the
// front page to ask -- that is the cycle this pair of constants exists to
// avoid -- so nothing but this catches them drifting, and a post would link
// back to a directory that is not there.
func TestPathMatchesTheListing(t *testing.T) {
	d := portfolio.DirByTab(portfolio.TabBlogs)

	if got := d.Path(); got != content.Path {
		t.Errorf("the listing says %q, a post page says %q", got, content.Path)
	}
	if got := d.Href(); got != content.Href {
		t.Errorf("the listing serves blogs at %q, a post page links to %q", got, content.Href)
	}
}

// TestFilterLinksOpenTheProjectsPage covers the one link that carries state:
// a tech chip on the stack page is a link to the projects page already
// filtered, since the signal that used to do it does not survive a navigation.
func TestFilterLinksOpenTheProjectsPage(t *testing.T) {
	stack := render(t, portfolio.TabStack)
	projects := portfolio.DirByTab(portfolio.TabProjects).Href()

	if !strings.Contains(stack, `href="`+projects+`?filter=`) {
		t.Errorf("no tech on the stack page links to a filtered %s", projects)
	}

	var b strings.Builder
	if err := PortfolioPage(portfolio.TabProjects, "Go").Render(context.Background(), &b); err != nil {
		t.Fatalf("rendering the projects page: %v", err)
	}
	if !strings.Contains(html.UnescapeString(b.String()), `"filter":"Go"`) {
		t.Error("a filter in the URL does not reach the page's signals")
	}
}

// TestEmptyDirectoriesStillSaySo covers the blogs shelf as it stands today: a
// directory with nothing in it is still a directory, and a page that rendered
// nothing at all for it would read as a bug rather than as an honest answer.
func TestEmptyDirectoriesStillSaySo(t *testing.T) {
	if len(content.Posts()) > 0 {
		t.Skip("there are posts now, so the empty shelf is not what renders")
	}
	if !strings.Contains(render(t, portfolio.TabBlogs), "Nothing written yet") {
		t.Error("the blogs page renders nothing at all when there are no posts")
	}
}
