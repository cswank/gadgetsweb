FROM golang

RUN sed -i 's/main$/main universe/' /etc/apt/sources.list &&\
    apt-get update &&\
    apt-get install libzmq3-dev &&\
    mkdir -p /opt/gadgets/bin /opt/gadgets/src/github.com/cswank/gadgetsweb /opt/gadgets/pkg

ADD . /opt/gadgets/src/bitbucket.org/cswank/gadgetsweb

RUN cd /opt/gadgets/src/bitbucket.org/cswank/gadgetsweb &&\
    export GOPATH=/opt/gadgets
    go get &&\
    go install

CMD = /opt/gadgets/bin/gadgetsweb
