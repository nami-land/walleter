FROM golang:latest

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
EXPOSE 8080
EXPOSE 8081

RUN go build
ENTRYPOINT ["./neco-wallet-center"]