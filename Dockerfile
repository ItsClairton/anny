FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target="/root/.cache/go-build" \
    go build -v -o bin .

FROM alpine:3.19

RUN --mount=type=cache,target=/etc/apk/cache apk add --update-cache ca-certificates \
    ffmpeg

WORKDIR /app

COPY --from=builder /app/bin /app/bin

CMD ["./bin"]