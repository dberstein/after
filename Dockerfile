FROM golang:1.20-alpine
RUN apk add --no-cache ncurses=6.3_p20221119-r0 make=4.3-r1 binutils=2.39-r2 alpine-sdk=1.0-r1 \
 && rm -vrf /var/cache/apk/*
WORKDIR /app
COPY . .
RUN make build