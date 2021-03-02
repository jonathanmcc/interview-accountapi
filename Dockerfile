FROM golang:1.14.15

ADD ./accountapi-client/ /go/src/

EXPOSE 8080

WORKDIR /go/src/

CMD ["go", "test"]