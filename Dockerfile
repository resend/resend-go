FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go install github.com/rakyll/gotest@latest

# RUN CGO_ENABLED=0 GOOS=linux go build -o /resend

# CMD ["/resend"]