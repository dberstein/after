FROM golang:1.20-alpine
RUN apk add make binutils
WORKDIR app
