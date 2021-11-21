
FROM golang:1.17.1-alpine 
EXPOSE 8888

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH  

RUN mkdir -p $GOPATH/src/dinning-hall
ADD . $GOPATH/src/dinning-hall

WORKDIR $GOPATH/src/dinning-hall 
RUN go build -o app . 

CMD ["/go/src/dinning-hall/app"]