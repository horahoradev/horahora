FROM golang:1.13.7-buster

COPY . /userservice

WORKDIR /userservice

RUN go build --mod=vendor -o userservice .

ENTRYPOINT ["./userservice"]
