FROM golang:1.10 as builder

RUN mkdir -p /go/src/github.com/f4hrenh9it/parley
WORKDIR /go/src/github.com/f4hrenh9it/parley
COPY . .
RUN CGO_ENABLED=0 go build -o target/bot ./cmd/bot

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/f4hrenh9it/parley/target/bot /bot
ENTRYPOINT ["./bot"]
