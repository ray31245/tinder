FROM golang:1.21 as builder
WORKDIR /workdir

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o tinder ./main.go

FROM alpine as base

RUN mkdir /workdir
WORKDIR /workdir

RUN apk --no-cache add libc6-compat

FROM base

COPY --from=builder /workdir/tinder /workdir/tinder

CMD [ "./tinder" ]