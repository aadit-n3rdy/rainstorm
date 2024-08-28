FROM golang:1.22-alpine

WORKDIR /go/src/

# RUN apk update && apk add git

RUN go install github.com/aadit-n3rdy/rainstorm/peer@628996702bcf0177548ff33ae2709ac046352a81
RUN go install github.com/aadit-n3rdy/rainstorm/tracker@628996702bcf0177548ff33ae2709ac046352a81

CMD [ "/bin/sh" ]
