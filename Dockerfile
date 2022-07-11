FROM golang:alpine as build
COPY src /src/
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /bin/server
RUN rm -rf /src

FROM golang:1.17.6-alpine as gcsfusebuilder
RUN apk add git
ARG GCSFUSE_REPO="/run/gcsfuse/"
ADD . ${GCSFUSE_REPO}
WORKDIR ${GCSFUSE_REPO}
RUN go install ./tools/build_gcsfuse
RUN build_gcsfuse . /tmp $(git log -1 --format=format:"%H")

FROM alpine:3.16.0
COPY --from=build /bin/server /bin/server
COPY --from=gcsfusebuilder /tmp/bin/gcsfuse /usr/local/bin/gcsfuse
COPY --from=gcsfusebuilder /tmp/sbin/mount.gcsfuse /usr/sbin/mount.gcsfuse
RUN apk add --no-cache tini fuse go
EXPOSE 8080
COPY gcsfuse_run.sh /app/gcsfuse_run.sh
RUN chmod +x /app/gcsfuse_run.sh
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/gcsfuse_run.sh"]