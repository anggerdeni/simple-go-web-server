FROM golang:alpine as build

COPY main.go go.* /src/
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /bin/server
RUN rm -rf /src


FROM alpine:3.16.0
COPY --from=build /bin/server /bin/server

ARG user=runner
ARG home=/home/$user
ENV PATH="/usr/local/bin:${PATH}"

RUN apk add --no-cache git make musl-dev go nodejs-current npm sudo
RUN echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home $home \
    --ingroup wheel \
    $user

RUN echo $user:pw | chpasswd
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
