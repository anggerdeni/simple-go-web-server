FROM golang:alpine as build

COPY main.go go.* /src/
WORKDIR /src
RUN CGO_ENABLED=0 go build -o /bin/server


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

RUN npm install -g jest@27.5.1 "@playwright/test@1.20.2" axios@0.27.1
RUN git clone https://github.com/anggerdeni/sample-rg-playground.git /opt/playground
EXPOSE 8080
ENTRYPOINT ["/bin/server"]
