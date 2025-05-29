FROM golang:1.24.3-alpine

WORKDIR /app
COPY . .
RUN go mod download

EXPOSE 8080

CMD ["go", "run", "cmd/api/main.go"]