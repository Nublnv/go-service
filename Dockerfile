# go-service
FROM golang:1.26.0 AS go-service
WORKDIR /src
COPY go-service/go.mod go-service/go.sum ./
RUN go mod download

COPY go-service/ .
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app ./cmd/app

RUN rm -rf /src/cmd
RUN groupadd -g 1000 nonroot && \
    useradd -m -u 1000 -g nonroot nonroot

EXPOSE 443

USER nonroot:nonroot
ENTRYPOINT [ "/app" ]

# py-service