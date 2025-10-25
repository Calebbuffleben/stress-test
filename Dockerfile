# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /loadtest ./cmd/loadtest

FROM alpine:3.20
LABEL name="stress-test"
RUN adduser -D -H appuser
WORKDIR /home/appuser
COPY --from=builder /server ./server
COPY --from=builder /loadtest ./loadtest
COPY entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh
ENV PORT=8080
EXPOSE 8080
USER appuser
ENTRYPOINT ["./entrypoint.sh"]


