FROM golang:1.15.2-buster

COPY . /frontend

WORKDIR /frontend

<<<<<<< HEAD
RUN apt-get update && \
    apt-get -y install nodejs npm && \
    npm install -g browserify && \
    cd javascript && npm install && make build

RUN go mod vendor && \
    go build --mod=vendor -gcflags "all=-N -l" -o frontend .
=======
RUN go build -gcflags "all=-N -l" -o frontend .
>>>>>>> 2e5526948c2f3902fc465d119f931cf419e3c3eb

ENTRYPOINT ["./frontend"]