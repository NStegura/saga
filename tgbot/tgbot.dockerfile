# Этап, на котором выполняется сборка приложения
FROM golang:1.22.0-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main cmd/tgbot/main.go

# Финальный этап, копируем собранное приложение
FROM alpine:3
COPY --from=builder /main /bin/main
ENTRYPOINT ["/bin/main"]