FROM golang:1.15.8-buster

COPY . /videoservice
WORKDIR /videoservice

#
#RUN go get github.com/go-delve/delve/cmd/dlv

# MP4box
RUN apt-get update && \
    apt-get install -y build-essential pkg-config git && \
    apt-get install -y zlib1g-dev && \
    git clone https://github.com/gpac/gpac gpac_public && \
    cd gpac_public && \
    ./configure --static-mp4box && \
    make -j8 && make install
    
RUN apt-get install -y ffmpeg curl

RUN go mod vendor && \
    go build --mod=vendor -o videoservice .

ENTRYPOINT ["./videoservice"]
