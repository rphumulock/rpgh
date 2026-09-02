package pages

import (
	"context"
	"html"
	"strings"
	"testing"

	"rpgh/features/blog/content"
	portfolio "rpgh/features/portfolio/components"
)

// render is the whole page as a visitor receives it. The panels are all in the
// document at once and hidden by $tab, so what a click does is decidable from
// the markup alone -- which is the only place the wiring exists.
func render(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	if err := PortfolioPage(portfolio.TabHome).Render(context.Background(), &b); err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	return html.UnescapeString(b.String())
}

// TestEveryDirectoryRowOpensAPanel is the guard the tab bar used to be. Nothing
// declares the set of panels twice any more: the listing sets $tab and a panel
// watches it, and a row naming a key no panel shows is a click that blanks the
// page with nothing in the console to say why.
func TestEveryDirectoryRowOpensAPanel(t *testing.T) {
	page := render(t)

	for _, d := range portfolio.Dirs() {
		if !strings.Contains(page, portfolio.SetTabExpr(d.Key)) {
			t.Errorf("nothing on the page opens %q", d.Name)
		}
		if !strings.Contains(page, portfolio.TabSelectedExpr(d.Key)) {
			t.Errorf("directory %q opens tab %q, which no panel shows itself for", d.Name, d.Key)
		}
	}
}

// TestEveryPanelCanGetBack checks the other direction now that there is no tab
// bar to fall back on: every directory carries a ../ home, or it is a panel a
// visitor can reach and then only leave by reloading.
func TestEveryPanelCanGetBack(t *testing.T) {
	back := strings.Count(render(t), portfolio.SetTabExpr(portfolio.TabHome))
	if want := len(portfolio.Dirs()); back != want {
		t.Errorf("the page has %d ways back to %s, want one per directory (%d)", back, portfolio.Root, want)
	}
}

// TestPanelsNameWhereTheyAre keeps the path bar honest: it is the only thing
// telling a visitor which directory they are standing in.
func TestPanelsNameWhereTheyAre(t *testing.T) {
	page := render(t)

	for _, d := range portfolio.Dirs() {
		if !strings.Contains(page, d.Path()) {
			t.Errorf("no panel names the path %q", d.Path())
		}
	}
}

// TestPageLandsOnTheListing guards the front page being the front page: the
// signal defaults are what decide which panel a visitor arrives at.
func TestPageLandsOnTheListing(t *testing.T) {
	if got := openOn(portfolio.TabHome).Tab; got != portfolio.TabHome {
		t.Errorf("the page opens on tab %q rather than the directory listing", got)
	}
}

// TestPathMatchesTheListing keeps the two spellings of ~/rpgh/blogs equal. The
// blog cannot import the front page to ask -- that is the cycle this pair of
// constants exists to avoid -- so nothing but this catches them drifting, and
// a post would say it came from a directory that is not on the listing.
func TestPathMatchesTheListing(t *testing.T) {
	if got := portfolio.DirByTab(portfolio.TabBlogs).Path(); got != content.Path {
		t.Errorf("the listing says %q, a post page says %q", got, content.Path)
	}
}

// TestEmptyDirectoriesStillSaySo covers the blogs shelf as it stands today: a
// directory with nothing in it is still a directory, and a panel that rendered
// an empty page for it would read as a bug rather than as an honest answer.
func TestEmptyDirectoriesStillSaySo(t *testing.T) {
	if len(content.Posts()) > 0 {
		t.Skip("there are posts now, so the empty shelf is not what renders")
	}
	if !strings.Contains(render(t), "Nothing written yet") {
		t.Error("the blogs panel renders nothing at all when there are no posts")
	}
}
