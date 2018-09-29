FROM golang:1.10.2-alpine

RUN apk add --no-cache --update alpine-sdk
COPY . /go/src/fileserver
RUN cd /go/src/fileserver && go build fileserver
FROM alpine:3.8

RUN mkdir /fileserver/
COPY --from=0 /go/src/fileserver/fileserver /fileserver/
WORKDIR /fileserver/

ENTRYPOINT ["./fileserver"]

CMD ["--port=8080"]
