FROM golang:latest as build
RUN mkdir -p /go/src/github.com/sunet/s3-mm-tool
ADD . /go/src/github.com/sunet/s3-mm-tool/
WORKDIR /go/src/github.com/sunet/s3-mm-tool
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN make
RUN env GOBIN=/usr/bin go install ./cmd/s3-mm-tool

# Now copy it into our base image.
FROM gcr.io/distroless/base:debug
COPY --from=build /usr/bin/s3-mm-tool /usr/bin/s3-mm-tool

ENTRYPOINT ["/usr/bin/s3-mm-tool","-s"]
