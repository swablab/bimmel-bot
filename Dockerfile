############################
# STEP 1 build executable binary
############################
FROM golang:1.17 AS builder 
RUN mkdir -p /swablab-bot
WORKDIR /swablab-bot
ADD . /swablab-bot
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app

############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /swablab-bot/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app"]