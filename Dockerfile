FROM golang:1.17

WORKDIR $GOPATH/src/github.com/amitbet/vncproxy

RUN mkdir -p $GOPATH/src/github.com/amitbet/vncproxy
COPY . .

RUN cd $GOPATH/src/github.com/amitbet/vncproxy/recorder/cmd && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /recorder .
RUN cd $GOPATH/src/github.com/amitbet/vncproxy/proxy/cmd && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -ldflags="-s -w" -o /proxy .
RUN cd $GOPATH/src/github.com/amitbet/vncproxy/player/cmd && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /player .

FROM scratch
COPY --from=0 /recorder /recorder
COPY --from=0 /proxy /proxy
COPY --from=0 /player /player

EXPOSE 5900

ENTRYPOINT ["/proxy"]