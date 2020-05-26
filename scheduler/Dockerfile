FROM golang:1.13.7-buster

COPY . /scheduler

WORKDIR /scheduler

# Fetch youtube-dl
RUN git clone https://github.com/animelover1984/youtube-dl.git  && \
    cd youtube-dl && \
    apt-get update && \
    apt-get install zip unzip && \
    apt-get install -y python3 python3-pip && \
    pip3 install -r /scheduler/youtube-dl/requirements.txt

RUN go build --mod=vendor -o scheduler .

ENTRYPOINT ["./scheduler"]
