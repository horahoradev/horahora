FROM golang:1.15.7-buster

RUN mkdir -p /scheduler/videos
WORKDIR /scheduler

RUN apt install -y git
RUN go get github.com/go-delve/delve/cmd/dlv

# Fetch youtube-dl
RUN git clone https://github.com/animelover1984/youtube-dl.git  && \
    cd youtube-dl && \
    apt-get update && \
    apt-get install zip unzip && \
    apt-get install -y python3 python3-pip && \
    pip3 install mutagen requests websockets python-dateutil && \
    apt-get install -y ffmpeg && \
    python3 /scheduler/youtube-dl/setup.py install

<<<<<<< HEAD
run git clone https://github.com/yt-dlp/yt-dlp.git && \
    cd yt-dlp && \
    pip3 install -r requirements.txt

RUN go mod vendor && \
    go build --mod=vendor -gcflags "all=-N -l" -o scheduler .
=======
COPY . /scheduler

RUN go build -gcflags "all=-N -l" -o scheduler .
>>>>>>> 2e5526948c2f3902fc465d119f931cf419e3c3eb

ENTRYPOINT ["./scheduler"]