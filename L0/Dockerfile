FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o /app/server -ldflags "-X=main.buildCommit=local -X=main.buildDate=today" ./cmd/main.go

FROM gcr.io/distroless/static:latest-arm64

COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]