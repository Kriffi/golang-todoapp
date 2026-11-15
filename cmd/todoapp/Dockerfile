FROM golang:1.27.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/cmd/todoapp/exe /app/cmd/todoapp


FROM migrate/migrate:v4.19.1 AS migrations

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/cmd/todoapp/exe /app

COPY --from=builder /app/public /app/public
COPY --from=migrations /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations
ENV PROJECT_ROOT=/app

CMD ["/app/exe"]
