FROM golang:1.24 as prod

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

CMD ["./app"]
