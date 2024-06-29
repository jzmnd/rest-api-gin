FROM golang:1.22 as builder

WORKDIR /go/src/rest-api-gin
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /main

EXPOSE 8080
CMD ["/main"]
