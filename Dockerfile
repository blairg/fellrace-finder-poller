# build stage
FROM golang:1.10.3-alpine3.8 AS build-env

WORKDIR /go/src/github.com/blairg/fellrace-finder-poller/

COPY ./ .

RUN apk --no-cache add git bzr mercurial && \
    go get -u github.com/golang/dep/... && \
    dep ensure -v --vendor-only && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# -------------------------------------------------------------------------------
# final stage
FROM alpine:latest  

ARG MONGO_DB_URL
ENV MONGO_DB_URL ${MONGO_DB_URL}
ENV RESULTS_PAGE_URL=http://www.fellrunner.org.uk/results.php

WORKDIR /root/

COPY --from=build-env /go/src/github.com/blairg/fellrace-finder-poller/app .

RUN apk --no-cache add ca-certificates

ENTRYPOINT ./app