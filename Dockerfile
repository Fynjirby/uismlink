FROM golang:1.24.4-alpine

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

ARG REPO_URL=https://github.com/fynjirby/uismlink.git

RUN git clone $REPO_URL . && \
    go mod tidy && \
    go build -o uismlink .

EXPOSE 8080
CMD ["./uismlink"]
