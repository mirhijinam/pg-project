FROM golang:1.22-alpine

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/app ./cmd/app/main.go

EXPOSE 7070

CMD ["./bin/app"]
