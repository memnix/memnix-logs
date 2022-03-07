FROM golang:1.17

RUN mkdir -p /go/src/memnix-logs
WORKDIR /go/src/memnix-logs

COPY . /go/src/memnix-logs

RUN go get -d -v
RUN go install -v

CMD ["/go/bin/memnixlogs"]