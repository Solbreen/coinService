FROM golang:1.19-alpine
WORKDIR /coinService
COPY . .
WORKDIR /coinService/cmd
RUN go mod download
RUN go build -o main .
CMD ["./main"]