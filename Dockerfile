FROM golang:1.14-alpine3.11 as builder
WORKDIR app
COPY app.go go.mod go.sum ./
RUN go build -o /app .

FROM alpine:3.11
EXPOSE 9000
CMD ["./app"]
COPY --from=builder /app .
