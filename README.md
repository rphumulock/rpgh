# rpgh

Personal website — a portfolio page and a resume, served by a single Go binary
with everything embedded.

# Stack

- [Go](https://go.dev/doc/)
- [Datastar](https://github.com/starfederation/datastar)
- [Templ](https://templ.guide/)
  - [Tailwind](https://tailwindcss.com/) x [DaisyUI](https://daisyui.com/)

# Pages

| Route     | Feature              | What it is                                                    |
| --------- | -------------------- | ------------------------------------------------------------- |
| `/`       | `features/portfolio` | Terminal theme, social links, project grid, stack tree         |
| `/resume` | `features/resume`    | The resume, plus the PDF it was transcribed from              |

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

# References

## Server

- [go](https://go.dev/)
- [datastar sdk](https://github.com/starfederation/datastar/tree/develop/sdk)
- [templ](https://templ.guide/)

## Client

- [datastar](https://www.jsdelivr.com/package/gh/starfederation/datastar)
- [tailwindcss](https://tailwindcss.com/)
- [daisyui](https://daisyui.com/)
