FROM golang:1.26.0 
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY go-service/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /app ./cmd/app

RUN rm -rf /src/cmd
RUN groupadd -g 1000 nonroot && \
    useradd -m -u 1000 -g nonroot nonroot

EXPOSE 443
ENV TLS_CERT=/src/certs/server.crt
ENV TLS_KEY=/src/certs/server.key

USER nonroot:nonroot
ENTRYPOINT [ "/app" ]
