FROM golang:1.24.1 AS builder

COPY ./src /app/src

WORKDIR /app/src

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vnc-summarizer


FROM gcr.io/distroless/static

COPY --from=builder /app/src/vnc-summarizer /vnc-summarizer

ENTRYPOINT ["/vnc-summarizer"]
