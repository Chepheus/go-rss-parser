# syntax=docker/dockerfile:1

FROM golang:1.22-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/web-parser

FROM alpine:3.20 as prod

WORKDIR /app

COPY --from=builder /app/web-parser .

EXPOSE 8080

RUN ln -s /app/web-parser /usr/local/bin/web-parser

CMD ["web-parser"]