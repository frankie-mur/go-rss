FROM golang:1.21.0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/web

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]