FROM golang:1.26-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o lab-monitor-backend .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /app/lab-monitor-backend .
# godotenv.Load() bikin app langsung Fatal kalau file .env ga ada,
# jadi filenya perlu ikut masuk image (nilainya nanti dioverride oleh
# environment: di docker-compose, godotenv ga nimpa env var yang udah diset).
COPY .env .

EXPOSE 8080
ENTRYPOINT ["./lab-monitor-backend"]
