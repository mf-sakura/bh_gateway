FROM golang:1.12.6-alpine3.10 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/mf-sakura/bh_gateway/app
COPY . /go/src/github.com/mf-sakura/bh_gateway/app

RUN apk update \
    && apk add --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure


RUN go build -o bh_gateway /go/src/github.com/mf-sakura/bh_gateway/app/cmd

# runtime image
FROM alpine
COPY --from=builder /go/src/github.com/mf-sakura/bh_gateway/app/bh_gateway /app/bh_gateway
WORKDIR /app
EXPOSE 5002
ENTRYPOINT ["/app/bh_gateway"]