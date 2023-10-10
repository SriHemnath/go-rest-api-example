FROM alpine:3.17 as root-certs
RUN apk add -U --no-cache ca-certificates
RUN mkdir /home/app
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.20 as builder
WORKDIR /go-rest-api-files
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o go-rest-api

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=1001:1001 --from=builder /go-rest-api-files/go-rest-api /go-rest-api
USER app
ENTRYPOINT ["./go-rest-api"]