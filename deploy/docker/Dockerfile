# Builder
FROM golang:1.23-alpine3.21 as builder

WORKDIR /usr/app

RUN apk update && apk upgrade && \
    apk --update add git make

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN make service

# Distribution
FROM alpine:3.21

WORKDIR /usr/app

COPY --from=builder /usr/app/bin/service /usr/app/bin/
COPY --from=builder /usr/app/config/ /usr/app/config/

ENV PATH="/usr/app/bin:${PATH}"

ENTRYPOINT [ "service" ]
