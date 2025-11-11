FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod init pod-watcher && go mod tidy && go build -o pod-watcher main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/pod-watcher .
USER nonroot:nonroot
ENTRYPOINT ["./pod-watcher"]
