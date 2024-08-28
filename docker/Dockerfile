FROM golang:1.22-alpine


WORKDIR /usr/src/

RUN go install github.com/aadit-n3rdy/rainstorm/peer@latest
RUN go install github.com/aadit-n3rdy/raindstorm/common@latest
RUN go install github.com/aadit-n3rdy/raindstorm/tracker@latest


CMD [ "/bin/sh" ]
