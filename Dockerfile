FROM golang

RUN apt-get update &&\
    apt-get -y install libzmq3-dev &&\
    mkdir -p /opt/gadgets/bin /opt/gadgets/src/github.com/cswank/gadgetsweb /opt/gadgets/pkg

ADD . /opt/gadgets/src/github.com/cswank/gadgetsweb

RUN GOPATH=/opt/gadgets go get github.com/vaughan0/go-zmq

RUN cd /opt/gadgets/src/github.com/cswank/gadgetsweb &&\
    export GOPATH=/opt/gadgets &&\
    go get &&\ 
    go install

CMD /opt/gadgets/bin/gadgetsweb
