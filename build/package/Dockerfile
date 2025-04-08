FROM golang:1.24 as prod

WORKDIR /app

COPY . /app

RUN go mod download

CMD ["go", "run", "cmd/main.go"]

FROM golang:1.24 as dev

WORKDIR /app

# Install air for live reloading (optional, remove if not using)
RUN go install github.com/air-verse/air@latest

# Copy only the go.mod and go.sum for dependency caching
COPY go.mod /app

# Download dependencies
RUN go mod download

WORKDIR /app/cmd

# live reload using `air`
CMD ["air"]