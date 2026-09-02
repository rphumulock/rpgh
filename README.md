# rpgh

Personal website — a portfolio page and a resume, served by a single Go binary
with everything embedded.

# Stack

- [Go](https://go.dev/doc/)
- [Datastar](https://github.com/starfederation/datastar)
- [Templ](https://templ.guide/)
  - [Tailwind](https://tailwindcss.com/) x [DaisyUI](https://daisyui.com/)

# Pages

| Route          | Feature              | What it is                                                |
| -------------- | -------------------- | --------------------------------------------------------- |
| `/`            | `features/portfolio` | A directory listing into projects, tech, blogs and videos  |
| `/projects`    | `features/portfolio` | Things built end to end, filterable by `?filter=<tech>`    |
| `/tech`        | `features/portfolio` | The toolbox as a tree                                      |
| `/blogs`       | `features/portfolio` | The posts, listed                                          |
| `/videos`      | `features/portfolio` | Recorded walkthroughs, grouped by series                   |
| `/blog/<slug>` | `features/blog`      | One post, rendered from the Markdown it is written in      |
| `/resume`      | `features/resume`    | The resume, plus the PDF it was transcribed from           |

Every directory on the front page is its own page at its own URL, so cd-ing
into one is a real navigation — the back button leaves it, and the address is a
link someone else can open on the same view. The listing rows and the `../` in
a panel's path bar are plain anchors, and the four directory routes are
registered from the same `Dirs()` the listing renders, so adding a directory
there serves it. The older `/?cd=blogs` links redirect to `/blogs`.

# Writing a post

Posts are Markdown files under
[features/blog/content/posts](./features/blog/content/posts), embedded into the
binary like everything else — publishing is a commit. The filename is the slug,
so `why-sse.md` is served at `/blog/why-sse`.

```markdown
---
title: Why SSE
published: 2026-09-02
blurb: One line, shown on the listing.
draft: true
---

The post itself.
```

Those four keys are the only ones recognised; anything else is an error rather
than a field that quietly does nothing. A post that does not parse fails the
tests and refuses to boot, with the filename in the message. `draft: true`
keeps one out of the listing and off its URL until you drop the line.

Each feature owns its own `routes.go`, `handlers.go`, and `pages/`; they are
registered together in [router/router.go](./router/router.go).

# Setup

Install dependencies:

```shell
go mod tidy
```

# Development

Live reload is set up out of the box — powered by
[Air](https://github.com/air-verse/air).

Use the [live task](./Taskfile.yml#L76) from the
[Taskfile](https://taskfile.dev/) to start with live reload:

```shell
go tool task live
```

Navigate to [`http://localhost:8080`](http://localhost:8080) to begin.

Templates are `.templ` files; after editing one, `templ generate` regenerates
the matching `_templ.go`. The `live` task does this for you on save.

## Debugging

The [debug task](./Taskfile.yml#L42) launches
[delve](https://github.com/go-delve/delve) against the project binary:

```shell
go tool task debug
```

## IDE Support

- [Templ / TailwindCSS Support](https://templ.guide/developer-tools/ide-support/)

### Visual Studio Code Integration

[Reference](https://code.visualstudio.com/docs/languages/go)

- [launch.json](./.vscode/launch.json)
- [settings.json](./.vscode/settings.json)

A `Debug Main` configuration has been added to the
[launch.json](./.vscode/launch.json).

# Starting the Server

```shell
go tool task run
```

Navigate to [`http://localhost:8080`](http://localhost:8080).

# Deployment

## Building an Executable

The `task build` [task](./Taskfile.yml#L33) will assemble and build a binary.

Static assets are served from disk under the `dev` build tag and embedded into
the binary otherwise — see
[web/resources](./web/resources/static_prod.go).

## Docker

```shell
# build an image
docker build -t rpgh:latest .

# run the image in a container
docker run --name rpgh -p 8080:9001 rpgh:latest
```

[Dockerfile](./Dockerfile)

## Continuous builds

[build.yml](./.github/workflows/build.yml) runs `go vet` and the tests on every
push and pull request, and fails if any committed `_templ.go` has drifted from
the `.templ` it came from. Pushes to `master` then publish the image to
`ghcr.io/rphumulock/rpgh`, tagged `latest` and with the commit sha alongside it
to roll back to.

The package is private until you make it public in the repository's Packages
settings; until then a pull from the server is denied.

## Running it

[compose.yaml](./compose.yaml) runs the site next to a
[Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/),
which dials out to Cloudflare rather than listening for anything — so no port
is published and no port is forwarded.

```shell
printf 'TUNNEL_TOKEN=...\n' > .env && chmod 600 .env
docker compose up -d
```

Updating is `docker compose pull && docker compose up -d`.

`TRUST_PROXY=true` tells the server to read the client address from
`CLIENT_IP_HEADER` (`CF-Connecting-IP`) rather than from the connection. Set it
only where something in front actually overwrites that header: it is what the
rate limiter keys on, and a header a proxy *appends* to — `X-Forwarded-For`,
notably — is the client's to forge.

# References

## Server

- [go](https://go.dev/)
- [datastar sdk](https://github.com/starfederation/datastar/tree/develop/sdk)
- [templ](https://templ.guide/)

## Client

- [datastar](https://www.jsdelivr.com/package/gh/starfederation/datastar)
- [tailwindcss](https://tailwindcss.com/)
- [daisyui](https://daisyui.com/)
