FROM golang:1.13.7-buster

COPY . /userservice

WORKDIR /userservice

# reference: https://blog.jetbrains.com/go/2018/04/30/debugging-containerized-go-applications/
RUN apt install -y git
RUN go get github.com/go-delve/delve/cmd/dlv

RUN go build --mod=vendor -gcflags "all=-N -l" -o userservice .

ENTRYPOINT ["./userservice"]
