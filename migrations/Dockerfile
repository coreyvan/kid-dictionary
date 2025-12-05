# Migration container using sql-migrate
# Built from source for multi-arch support (including ARM64)

FROM golang:1.25-alpine AS builder

RUN go install github.com/rubenv/sql-migrate/...@latest

FROM alpine:3.20

COPY --from=builder /go/bin/sql-migrate /usr/local/bin/sql-migrate

WORKDIR /app

ENTRYPOINT ["sql-migrate"]
CMD ["up", "-config=dbconfig.yml"]
