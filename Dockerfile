FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/

FROM alpine:latest
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
