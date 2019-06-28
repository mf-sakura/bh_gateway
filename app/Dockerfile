FROM golang:1.12.6-alpine3.10 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/mf-sakura/bh_user/app
COPY . /go/src/github.com/mf-sakura/bh_user/app

RUN apk update \
    && apk add --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure


RUN go build -o bh_user /go/src/github.com/mf-sakura/bh_user/app/cmd

# runtime image
FROM alpine
COPY --from=builder /go/src/github.com/mf-sakura/bh_user/app/bh_user /app/bh_user
WORKDIR /app
EXPOSE 5002
ENTRYPOINT ["/app/bh_user"]