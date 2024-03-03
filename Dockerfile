FROM golang:latest

WORKDIR /

ADD go.mod .

ADD go.sum .

RUN go mod download

RUN apt-get update && apt-get install sqlite3

COPY . .

RUN go build -o apiserver ./cmd/apiserver/

CMD ["./apiserver"]

