FROM golang:1.13.7-buster

COPY . /scheduler

WORKDIR /scheduler

RUN apt install -y git
RUN go get github.com/go-delve/delve/cmd/dlv

# Fetch youtube-dl
RUN git clone https://github.com/animelover1984/youtube-dl.git  && \
    cd youtube-dl && \
    apt-get update && \
    apt-get install zip unzip && \
    apt-get install -y python3 python3-pip && \
    pip3 install -r /scheduler/youtube-dl/requirements.txt

RUN go build --mod=vendor -gcflags "all=-N -l" -o scheduler .

ENTRYPOINT ["./scheduler"]
