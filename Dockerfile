# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Compila o binário estaticamente
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stress-test .
RUN chmod +x stress-test 

# Imagem final pequena
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/stress-test .
RUN chmod +x stress-test 
ENTRYPOINT ["./stress-test"]
