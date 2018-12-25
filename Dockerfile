FROM golang:1.11.2-alpine3.8
MAINTAINER M. - Karan Bhomia

ENV SOURCES /go/src/github.com/karanbhomiagit/reverse-proxy-aggregation-api/

COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go install

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT reverse-proxy-aggregation-api
