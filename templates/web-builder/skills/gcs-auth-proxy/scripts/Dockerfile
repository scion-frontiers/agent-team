FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 go build -o proxy .

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/proxy /proxy
ENTRYPOINT ["/proxy"]
