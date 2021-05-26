FROM golang:1.15.2-buster

COPY . /frontend

WORKDIR /frontend

RUN go build -gcflags "all=-N -l" -o frontend .

ENTRYPOINT ["./frontend"]