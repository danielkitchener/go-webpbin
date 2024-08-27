FROM golang

WORKDIR /

RUN go get -v github.com/belphemur/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v golang.org/x/image/webp

RUN mkdir -p $GOPATH/src/github.com/belphemur/go-webpbin
COPY . $GOPATH/src/github.com/belphemur/go-webpbin
WORKDIR $GOPATH/src/github.com/belphemur/go-webpbin
RUN go test -v ./...