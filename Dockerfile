FROM docker.io/golang:1.23-alpine3.19 AS builder
WORKDIR /app
COPY . .

RUN go mod download

RUN GODEBUG=asyncpreemptoff=1 go build -o /app/fren cmd/fren/fren.go
RUN GOBIN=/app go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.19.1
WORKDIR /app
COPY --from=builder /app/fren /app/fren
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/migrate /app/migrate

EXPOSE 8080

ENTRYPOINT ["/app/fren"]
