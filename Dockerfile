FROM ghcr.io/hybridgroup/opencv:4.10.0
ENV TZ=Asia/Tokyo
ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN go build -o /bin/bot ./pkg
ENTRYPOINT ["/bin/bot"]
