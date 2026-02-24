FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o load-generator .

FROM alpine:3.19

COPY --from=builder /app/load-generator /usr/local/bin/load-generator

EXPOSE 8080

ENTRYPOINT ["load-generator"]
