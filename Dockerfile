# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/app main.go

# Stage 2: Minimal runtime image
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
USER appuser
COPY --from=builder /bin/app /bin/app
EXPOSE 8000
ENTRYPOINT ["/bin/app"]
