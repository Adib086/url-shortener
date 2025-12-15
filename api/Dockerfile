FROM golang:alpine AS builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go build -o main .

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

WORKDIR /app

COPY --from=builder /build/main /app/

EXPOSE 3000

CMD ["./main"]

