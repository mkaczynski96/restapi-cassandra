FROM golang:latest
ADD . /go/src/acaisoft-mkaczynski-api
RUN go get github.com/gocql/gocql
RUN go get github.com/gorilla/mux
RUN go install acaisoft-mkaczynski-api
ENTRYPOINT /go/bin/acaisoft-mkaczynski-api
EXPOSE 8080