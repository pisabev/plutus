FROM golang:1.21.1

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/api ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./cmd/main.go

CMD ["/api/app"]