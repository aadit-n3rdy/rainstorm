FROM golang:1.22-alpine

WORKDIR /go/src/

# RUN apk update && apk add git

RUN go install github.com/aadit-n3rdy/rainstorm/peer@main
RUN go install github.com/aadit-n3rdy/rainstorm/tracker@main

CMD [ "/bin/sh" ]
