FROM docker.io/golang:alpine AS build

RUN apk add --no-cache upx

WORKDIR /src
COPY . ./
RUN go mod download

# Generate rather than trust what was copied in: the templ output is committed
# but can drift from the .templ sources, and index.css is gitignored entirely,
# so a clean checkout has no stylesheet to embed at all.
RUN go tool templ generate
RUN go tool gotailwind -i web/resources/styles/styles.css -o web/resources/static/index.css

RUN --mount=type=cache,target=/root/.cache/go-build \
go build -ldflags="-s -w" -o /bin/main ./cmd/web
RUN upx -9 -k /bin/main

FROM scratch
ENV PORT=9001
COPY --from=build /bin/main /
# Nobody. The binary reads nothing off disk and writes nothing, so it has no
# reason to be uid 0 even in an image with no shell to drop into.
USER 65534:65534
ENTRYPOINT ["/main"]
