FROM golang:1.22 as builder

WORKDIR /go/src/
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o app
CMD ["./app"]