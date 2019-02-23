# build gjo
FROM golang AS builder
ENV GOPATH /go
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
ENV GO111MODULE on
COPY . ./src/github.com/skanehira/gjo
WORKDIR /go/src/github.com/skanehira/gjo
RUN go build

# copy artifact from the build stage
FROM busybox
COPY --from=builder /go/src/github.com/skanehira/gjo/gjo /usr/local/bin/gjo

ENTRYPOINT ["gjo"]
