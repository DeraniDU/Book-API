# Use the official Go image to build the app
FROM golang:1.24.1 AS builder  

# Rest of the Dockerfile remains the same
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o book-api main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/book-api .
EXPOSE 8080
CMD ["./book-api"]