FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* .
RUN go mod download

COPY . .

RUN go build -o main ./cmd/bot/main.go

# Runtime stage
FROM alpine:latest

RUN touch .env

RUN echo CONFIG_PATH=${{ secrets.PATH }} >> .env
RUN echo CONFIG_PATH=${{ secrets.BATCH_SIZE }} >> .env
RUN echo CONFIG_PATH=${{ secrets.TELEGRAM_TOKEN }} >> .env
RUN echo CONFIG_PATH=${{ secrets.TELEGRAM_HOST }} >> .env
RUN echo CONFIG_PATH=${{ secrets.MAIL_PORT }} >> .env
RUN echo CONFIG_PATH=${{ secrets.MAIL_PASSWORD }} >> .env
RUN echo CONFIG_PATH=${{ secrets.MAIL_HOST }} >> .env
RUN echo CONFIG_PATH=${{ secrets.MAIL_FROM }} >> .env
RUN echo CONFIG_PATH=${{ secrets.MAIL_TO }} >> .env

COPY --from=builder /app/main main

ENTRYPOINT ["/main"]