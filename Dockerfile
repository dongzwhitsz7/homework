FROM golang:1.17 as builder
#ENV GO111MOULE=on CGO_ENABLE=0 GOOS=linux GOARCH=amd64
WORKDIR /build
COPY . .
RUN go build -o httpserver

#FROM scratch
#COPY --from=builder /build/httpserver /
EXPOSE 80
ENTRYPOINT ["./httpserver"]