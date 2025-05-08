FROM golang:1.24.0-alpine AS builder

WORKDIR /workspace

COPY ../.. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

RUN apk add --no-cache curl
RUN curl https://letsencrypt.org/certs/isrgrootx1.pem > ./ca-cert.pem

FROM scratch

COPY --from=builder /workspace/app /app
COPY --from=builder /workspace/ca-cert.pem /etc/ssl/certs/ca-cert.pem
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/app"]
