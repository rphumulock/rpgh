FROM docker.io/golang:alpine AS build

RUN apk add --no-cache upx

WORKDIR /src
COPY . ./
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -ldflags="-s -w" -o /bin/main ./cmd/web
RUN upx -9 -k /bin/main

FROM scratch
ENV PORT=9001
COPY --from=build /bin/main /
ENTRYPOINT ["/main"]
