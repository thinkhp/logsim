FROM golang:1.10

RUN mkdir -p go/src/github.com/thinkhp/logsim
WORKDIR go/src/github.com/thinkhp/logsim

RUN go get -v github.com/pkg/errors
COPY ./log_test.go .
COPY ./logger.go .
COPY ./log_rotate.go .

CMD ["go", "test", "-v"]