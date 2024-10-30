FROM golang:1.23 AS builder

ARG GITHUB_USER=$GITHUB_USER
ARG GITHUB_PASSWORD=$GITHUB_PASSWORD

RUN echo "machine github.com\n\tlogin $GITHUB_USER\n\tpassword $GITHUB_PASSWORD" >> ~/.netrc

COPY ./src /app/src

WORKDIR /app/src

RUN go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vnc-summarizer


FROM gcr.io/distroless/static

COPY --from=builder /app/src/vnc-summarizer /vnc-summarizer

ENTRYPOINT ["/vnc-summarizer"]
