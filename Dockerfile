# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
