FROM golang:1.21.5
WORKDIR /src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
CMD ["go", "run", "./cmd/web"]