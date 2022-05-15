FROM golang:1.18.2-alpine3.15 as builder
WORKDIR /go/src/app
COPY . /go/src/app
RUN apk add git
RUN go build -i main.go


FROM  alpine:3.15
WORKDIR /root/
COPY --from=builder /go/src/app/main .
CMD ["./main"]  