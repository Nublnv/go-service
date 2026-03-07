# go-service
FROM golang:1.26.0 AS build
WORKDIR /src
COPY go-service/go.mod go-service/go.sum ./
RUN go mod download

COPY go-service/ .
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -trimpath -ldflags="-s -w" -o /app ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -trimpath -ldflags="-s -w" -o /healthcheck ./healthcheck

FROM debian:stable-slim AS runtime

RUN groupadd -g 1000 nonroot && \
useradd -m -u 1000 -g nonroot nonroot

WORKDIR /

COPY --from=build --chown=nonroot:1000 /src/migrations /src/migrations
COPY --from=build --chown=nonroot:1000 /app /app
COPY --from=build --chown=nonroot:1000 /healthcheck /healthcheck

EXPOSE 443

USER nonroot:nonroot
ENTRYPOINT [ "/app" ]

# py-service