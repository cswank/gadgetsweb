FROM golang

RUN sed -i 's/main$/main universe/' /etc/apt/sources.list && apt-get update

# prevent apt from starting postgres right after the installation
RUN echo "#!/bin/sh\nexit 101" > /usr/sbin/policy-rc.d; chmod +x /usr/sbin/policy-rc.d

RUN apt-get install libzmq3-dev

RUN mkdir -p /opt/gadgets/bin
RUN mkdir -p /opt/gadgets/src/bitbucket.org/cswank/gadgetsweb
RUN mkdir /opt/gadgets/pkg

ADD . /opt/gadgets/src/bitbucket.org/cswank/gadgetsweb



