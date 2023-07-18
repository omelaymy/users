# Build
FROM golang:latest AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/users ./cmd/api/main.go

# Launch
FROM alpine:latest

EXPOSE 8888

WORKDIR /app

COPY --from=build /app/users /app/users
COPY --from=build /app/config /app/config

CMD ["sh", "-c", "./users"]