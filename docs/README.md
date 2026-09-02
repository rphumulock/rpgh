# Decisions

Why the site is shaped the way it is. The code says what it does and the commit
messages say what changed; this is the layer above both — the choices that are
still load-bearing, and the ones that will look arbitrary later without the
reason attached.

Start at [the root README](../README.md) for what the project is and how to run
it.

## The site is a filesystem

The landing page is `ls -l ~/rpgh`: four rows, one per directory, with a
permission column, a count where a byte size would go, the name, and a blurb.

That came from replacing a front page that opened straight onto the projects
grid with everything else behind tabs. The problem with tabs was not how they
looked — it was that they put the tab bar in charge of announcing what the site
holds, and a visitor landing on projects had to read a row of labels to learn
there was anything else. A directory listing says it in one glance, in the
idiom the rest of the site was already speaking.

**Counts are read, never written.** A row reports `len(Projects)`,
`StackCount()`, `len(content.Posts())`, `TotalVideoCount()`. Adding a project is
one data edit, not two, and the listing cannot drift from what the panel
renders.

## Every directory is a URL

`/`, `/projects`, `/tech`, `/blogs`, `/videos`.

The first version of this was one document holding every panel, switched by a
`$tab` signal. It worked, and it was wrong: cd-ing into a directory was not a
navigation, so the browser had nothing to go back to, and a directory had no
address anyone could send. The `../` in the path bar was the only way out of
one.

Routes are registered from the same `Dirs()` the listing renders, so a
directory added there is served without a second edit — and there is no way to
add one the listing links to that nothing serves.

**Key and route are separate on purpose.** `TabStack` is the Go identifier;
`/tech` is the address. `tech/` was spelled `stack` in the code before it was
ever a route, and renaming an identifier to match a URL is the kind of churn
that makes a diff unreadable for no gain.

### State that had to survive the change

One click used to carry state across panels: a tech chip on the stack page set
a filter and jumped to projects. A signal does not survive a navigation, so it
links to `/projects?filter=<tech>` and the handler seeds the signal from the
query. A tech nothing was built with falls back to showing everything — a URL
gets typed by hand sometimes, and an empty grid would read as a bug rather than
as an answer.

### Old links still resolve

While the site was one page, a directory was reachable as `/?cd=blogs`, and post
pages linked back that way. Those links are out in the world, so the root
redirects them (301) to the route the directory lives at now, rather than
quietly landing on the listing. Nothing in the source emits `?cd=` any more; the
redirect exists only for what already escaped.

## `../` instead of a tab bar

Each panel opens with a bar carrying `../` and the path you are standing in. It
is pinned above the panel's own scroll rather than inside it, so the way out
does not scroll away.

A tab bar and a directory listing pointing at the same four places is the same
navigation twice. Since a panel is a directory, the way out of one is the `../`
you would have typed.

## A wrong path gets the listing

chi's bare 404 was the one page on the site with none of the site on it.

The answer is the shell's — `cd: ~/rpgh/nope: No such file or directory` — with
the directory listing underneath, because the useful half of "no such
directory" is knowing which ones there are. That is why the handler lives with
the portfolio: the listing is the front page's, and a wrong turn should end on
the page that has somewhere to go.

`/blog/<slug>` keeps its own 404, which names the file it went looking for
(`cat: ~/rpgh/blogs/nope.md: No such file or directory`) rather than the
directory.

## Videos are a directory, not a footnote

This reverses [`dc76e97`](https://github.com/rphumulock/rpgh/commit/dc76e97),
which had moved videos under the project grid.

Once the front page became a listing, videos were one of the things worth
listing rather than a section appended to projects. The section now owns its own
scroll the way the stack panel does, and lost the top border that had separated
it from the grid above.

## The blog

Posts are Markdown under `features/blog/content/posts`, embedded in the binary.
Publishing is a commit. There is no CMS, no database, and nothing to back up
separately from the repository.

**The filename is the URL.** `why-sse.md` is served at `/blog/why-sse`. No slug
field to keep in sync with anything, and no table mapping one to the other.

**Frontmatter is four keys and nothing else** — `title`, `published`, `blurb`,
`draft`. An unrecognised key is an error, not a field that silently does
nothing. This is deliberately not a YAML parser: a real one would accept
frontmatter that the rest of the code has no idea what to do with, and the
failure would be a page that renders wrong rather than a build that stops.

**Parsing happens once, at boot.** `SetupRoutes` returns the error, so a
malformed post fails the tests in CI and refuses to start the server, with the
filename in the message. The alternative — parsing per request — turns an
author's typo into something a reader discovers.

**`draft: true` means both halves.** Kept off the listing *and* unreachable by
guessing the URL, so an unfinished post cannot be read by anyone who knows the
filename convention.

**An empty directory still prints.** With no posts, the panel says `total 0` and
shows an empty shelf. Hiding the section until there was something in it would
mean the front page counts a directory whose row opens onto a blank page, which
reads as broken rather than as an honest answer.

Rendering is [goldmark](https://github.com/yuin/goldmark) — the one dependency
this added. Its output has no classes to hang Tailwind on, so post bodies are
styled by element under a single `.prose-terminal` wrapper in `styles.css`,
scoped so it cannot reach the rest of the page.

## Package layout was forced by import cycles

Two arrangements that look natural do not compile, and the current shape is what
is left after both were ruled out.

**Site chrome lives in `features/common/components`.** The post page needs the
header and footer. Those were in `features/portfolio/pages`, and importing that
from `features/blog` cycles. Moving them also fixed a wart that predated the
blog: the resume was already reaching into the portfolio for chrome that is not
the portfolio's.

**Post data lives in `features/blog/content`, apart from the routes in
`features/blog`.** The portfolio needs the post count for the listing; the blog
needs the portfolio's chrome for its pages. With data and handlers in one
package that is `blog → blog/pages → portfolio/pages → portfolio/components →
blog`. Splitting the data out breaks it, because nothing imports
`features/blog`.

## Two places spell the same string

`~/rpgh/blogs` is written in `content.Path` and built again by the front page as
`Root + "/" + Name`. `content` cannot import the listing to ask — that is the
cycle above.

`TestPathMatchesTheListing` is the only thing holding them equal. Without it a
post page would name a directory that is not on the listing, and nothing else
would notice.

This is the repository's usual answer to a fact that has to exist twice — see
also `TestTechNamesResolve`, which keeps `Project.Tech` and `Stack` in sync.

## `.gitignore` is an allowlist, and it bites

The file ignores `*` and un-ignores specific paths. `*.md` is **not** among
them — only `README.md` is, which is why this document is `docs/README.md` and
not `docs/decisions.md`.

Two consequences that were nearly shipped broken:

- Blog posts needed `!/features/blog/content/posts/*.md`. Without it the first
  post would have been written, committed, and silently absent.
- The posts directory needs its `.gitkeep` tracked. `features/blog/content`
  embeds that directory, and a clone missing it does not compile — an `embed`
  pattern that matches nothing is a build error, which is also why the pattern
  is `all:posts` rather than `posts/*.md`.

**Before adding any non-Go, non-templ file to this repo, check
`git check-ignore -v <path>`.** A file that is silently not committed looks
exactly like a file that is.

## Content Security Policy

### The nonce buys one thing

128 bits of `crypto/rand`, minted per response in `SecurityHeaders`, set in the
header and stamped on `<html data-nonce>`. Datastar reads it at startup, removes
the attribute, and compiles every `data-*` expression into a script element
carrying it.

The choice is not "nonce or nothing" — it is **nonce or `'unsafe-eval'`**.
Without one, Datastar falls back to the `Function` constructor and the policy
has to allow eval outright, which reopens string-to-code for the whole page.
`policy()` still has that fallback branch, because blocking Datastar there
breaks every signal on the page rather than hardening it.

`'strict-dynamic'` is deliberately absent: it would make the host allowlist be
ignored, and iconify would stop loading.

**`Cache-Control: no-store` on pages follows from the nonce**, and is not a
separate decision. A page carrying one nonce, cached and re-served to someone
whose header names another, has scripts that do not match — and a nonce shared
between visitors is not unpredictable, so it is not a nonce. Static assets carry
none and cache hard.

### The analytics beacon is allowlisted by host

Cloudflare injects the Web Analytics beacon at the edge, *after* this server has
written the policy header. It can never carry our nonce, so
`https://static.cloudflareinsights.com` is in `script-src` by hostname or the
browser refuses it.

Under automatic injection the beacon reports to this origin's `/cdn-cgi/rum`,
which `connect-src 'self'` already covers. **Installing the snippet by hand
instead would post to `https://cloudflareinsights.com` and need that added to
`connect-src`** — a different fix, and the one thing to get right if the
approach ever changes.

### Why the policy has tests

A dropped source in `script-src` breaks no build and no test. The browser
refuses one script, the site otherwise works, and the only signal is a console
error.

For the beacon specifically, that error is indistinguishable from DNS
null-routing the host — which, on this network, is the first and usually correct
suspicion. `router/middleware_test.go` makes it decidable: **red means the CSP,
green means it is DNS.**

## Deliberately not here

- **No client-side router.** Every directory is a server route and a full page
  load. The pages are small and the server renders them in microseconds.
- **No per-post route table.** The filename is the slug; a post is a file.
- **No database.** Everything the site knows is a Go value or a Markdown file in
  this repository, and the binary carries all of it.
- **No `?tab=` on the root.** A directory has an address; the root is the
  listing and nothing else.
