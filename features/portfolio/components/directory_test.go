package components

import "testing"

// TestNotFoundPathReadsAsAPath keeps the 404's headline in the same voice as
// the path bar. It is printed straight into `cd ...`, so a request path that
// came back doubled up on slashes or missing the root would read as the site
// being confused rather than the visitor.
func TestNotFoundPathReadsAsAPath(t *testing.T) {
	for _, tc := range []struct{ url, want string }{
		{"/nope", Root + "/nope"},
		{"/nope/", Root + "/nope"},
		{"/a/b", Root + "/a/b"},
		{"nope", Root + "/nope"},
		{"/", Root},
		{"", Root},
	} {
		if got := NotFoundPath(tc.url); got != tc.want {
			t.Errorf("NotFoundPath(%q) = %q, want %q", tc.url, got, tc.want)
		}
	}
}

// TestDirectoriesAreServedWhereTheyAreListed pins the two spellings of a
// directory to each other: the name in the path bar and the name in the URL.
// They are the same word by design, so someone reading ~/rpgh/tech can type
// /tech and arrive.
func TestDirectoriesAreServedWhereTheyAreListed(t *testing.T) {
	for _, d := range Dirs() {
		if want := "/" + d.Name; d.Href() != want {
			t.Errorf("%q is listed as %q but served at %q", d.Key, want, d.Href())
		}
	}
}
