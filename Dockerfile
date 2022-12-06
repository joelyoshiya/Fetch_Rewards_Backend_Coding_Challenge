# alpine chosen for small footprint
FROM golang:1.18.3-alpine

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./main/* .
# compile the app
RUN go build -o /receipt-processor-service

EXPOSE 8080

CMD ["/receipt-processor-service"]