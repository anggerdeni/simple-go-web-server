FROM golang:alpine as build
COPY src /src/
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /bin/server
RUN rm -rf /src

FROM golang:alpine AS gcsfusebuilder
ENV GO111MODULE=off
ARG GCSFUSE_VERSION=0.41.4
RUN apk --update --no-cache add git fuse fuse-dev;
RUN go get -d github.com/googlecloudplatform/gcsfuse
RUN go install github.com/googlecloudplatform/gcsfuse/tools/build_gcsfuse
RUN build_gcsfuse ${GOPATH}/src/github.com/googlecloudplatform/gcsfuse /tmp ${GCSFUSE_VERSION}


FROM alpine:3.16.0
COPY --from=build /bin/server /bin/server
COPY --from=gcsfusebuilder /tmp/bin/gcsfuse /usr/bin
COPY --from=gcsfusebuilder /tmp/sbin/mount.gcsfuse /usr/sbin
RUN ln -s /usr/sbin/mount.gcsfuse /usr/sbin/mount.fuse.gcsfuse
RUN apk add --no-cache tini fuse go 
EXPOSE 8080
COPY gcsfuse_run.sh /app/gcsfuse_run.sh
RUN chmod +x /app/gcsfuse_run.sh
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/gcsfuse_run.sh"]