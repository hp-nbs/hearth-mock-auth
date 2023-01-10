FROM golang:1.18-alpine

WORKDIR /app/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /bin/auth ./cmd/main.go

EXPOSE 7000
ENTRYPOINT ["/bin/auth"]