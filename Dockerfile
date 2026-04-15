FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o meowcut ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/meowcut .

COPY configs ./configs
COPY migrations ./migrations

RUN mkdir -p logs

EXPOSE 8080

CMD ["./meowcut"]
