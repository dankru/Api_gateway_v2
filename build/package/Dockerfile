FROM golang:1.24 as prod

WORKDIR /app

COPY . /app

RUN go mod download

CMD ["go", "run", "cmd/main.go"]

FROM golang:1.24 as dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod /app

RUN go mod download

WORKDIR /app/cmd

# live reload using `air`
CMD ["air"]