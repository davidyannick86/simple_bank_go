# Build Stage
FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add --no-cache curl \
    && curl -L "https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz" | tar xvz \
    && chmod +x migrate

# Run Stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY db/migrations ./migrations
COPY start.sh .

RUN chmod +x start.sh

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]