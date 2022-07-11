FROM golang:alpine as build

COPY src /src/
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /bin/server
RUN rm -rf /src


FROM alpine:3.16.0
COPY --from=build /bin/server /bin/server

ARG user=runner
ARG home=/home/$user
ENV PATH="/usr/local/bin:${PATH}"

RUN apk add --no-cache go nodejs-current npm
EXPOSE 8080

# Use tini to manage zombie processes and signal forwarding
# https://github.com/krallin/tini
ENTRYPOINT ["/usr/bin/tini", "--"] 

# Pass the startup script as arguments to Tini
CMD ["/app/gcsfuse_run.sh"]