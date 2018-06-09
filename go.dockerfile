FROM golang:1.10.3

WORKDIR /go/src/app/ginserver
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

# RUN dep init
# dep creates a fresh clone of all the dependencies on $GOPATH/pkg/dep/sources/
# Gopkg.lock and Gopkg.toml are originated from the command
# use Gopkg.toml to specify different versions of a dependency

RUN dep ensure -update
