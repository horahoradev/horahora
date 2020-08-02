FROM golang:1.13.7-buster

COPY . /frontend

WORKDIR /frontend

RUN go build --mod=vendor -gcflags "all=-N -l" -o frontend .

ENTRYPOINT ["./frontend"]