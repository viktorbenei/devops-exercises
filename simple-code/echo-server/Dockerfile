FROM golang:1.11

# build dependencies
# RUN apt-get update -qq \
#  && DEBIAN_FRONTEND=noninteractive apt-get -y install \
#     openssl ca-certificates

# build the binary
ADD . /go/src/github.com/viktorbenei/devops-exercises/simple-code/echo-server
# The `-mod=vendor` flag instructs the go tool to use the vendor dir,
#  otherwhise the vendor/ dir is ignored in 1.11 module mode.
#  For more info see: https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
RUN cd /go/src/github.com/viktorbenei/devops-exercises/simple-code/echo-server \
    && go vet -v ./... \
    # && go get -u github.com/kisielk/errcheck && errcheck -asserts=true -blank=true -verbose ./... \
    # && go build -mod=vendor -o /go/bin/server ./server
    && go build -o /go/bin/server .

# FROM alpine:latest
FROM ubuntu:18.04
COPY --from=0 /go/bin/server /app/server

# runtime dependencies
RUN apt-get update -qq \
    && DEBIAN_FRONTEND=noninteractive apt-get -y install \
    ca-certificates

CMD ["/app/server"]