FROM golang:1.9.2

WORKDIR /go/src/github.com/someanon/golang-test-task
COPY . .

RUN go-wrapper download
RUN go-wrapper install

ENTRYPOINT golang-test-task

EXPOSE 8080