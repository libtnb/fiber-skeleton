FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=${VERSION}" -o /out/app ./cmd/app \
 && CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=${VERSION}" -o /out/cli ./cmd/cli

FROM alpine:3
RUN adduser -D -u 1000 app
WORKDIR /app
COPY --from=build /out/app /out/cli ./
RUN mkdir -p config storage/logs && chown -R app:app /app
USER app
EXPOSE 3000
VOLUME ["/app/config", "/app/storage"]
# adjust the port if you change http.address
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -q -O /dev/null http://127.0.0.1:3000/healthz || exit 1
ENTRYPOINT ["/app/app"]
