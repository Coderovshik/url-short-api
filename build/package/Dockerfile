FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 go build -o ./bin/url-short ./cmd/url-short/main.go

RUN --mount=type=cache,target=/go/pkg/ go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine AS runner

WORKDIR /

COPY --from=builder /go/bin/goose goose
COPY --from=builder /app/bin/url-short url-short

COPY  configs/* configs/
COPY sql/* sql/
COPY scripts/* scripts/

EXPOSE 8080

ENTRYPOINT [ "/scripts/startup.sh" ]