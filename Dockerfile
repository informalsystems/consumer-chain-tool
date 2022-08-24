FROM fedora:36
ENV GOPATH=/go
ENV PATH=$PATH:/go/bin
RUN dnf update -y
RUN dnf install -y golang gcc gcc-c++ jq procps

ADD . /consumer-chain-tool

RUN pushd /consumer-chain-tool/ && PATH=$PATH:/usr/local/go/bin GOPROXY=https://proxy.golang.org && PATH=$PATH:/usr/local/go/bin go install

# TODO: Once the binaries are final they will be downloaded from somewhere
COPY ./wasmd /go/bin/
COPY ./wasmd_consumer /go/bin/

WORKDIR /go/bin/