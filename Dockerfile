FROM golang:1.24 as prod

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o app ./cmd
CMD ["./app"]
