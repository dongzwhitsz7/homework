FROM golang:1.17 AS builder
ENV GO111MODULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY . .
RUN go build -o httpserver .

FROM scratch
COPY --from=builder /build/httpserver /
EXPOSE 8888
ENTRYPOINT ["/httpserver"]