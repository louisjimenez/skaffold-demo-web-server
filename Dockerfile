FROM golang:1.14-alpine3.11 as builder
COPY app.go .
RUN go build -o /app .

FROM alpine:3.11
CMD ["./app"]
COPY --from=builder /app .