FROM golang:1.21.5
WORKDIR /src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN mkdir bin
RUN GOOS=linux go build -o ./bin/web ./cmd/web


EXPOSE 4000

CMD ["./bin/web"]
