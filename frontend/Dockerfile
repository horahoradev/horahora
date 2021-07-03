FROM golang:1.15.2-buster

COPY . /frontend

WORKDIR /frontend

RUN apt-get update && \
    apt-get -y install nodejs npm && \
    npm install -g browserify && \
    cd javascript && npm install && make build

RUN go mod vendor && \
    go build --mod=vendor -gcflags "all=-N -l" -o frontend .

ENTRYPOINT ["./frontend"]