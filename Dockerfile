FROM golang:1.17

RUN apt-get update && \
    apt install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

RUN go env -w GO111MODULE="off" && \
    go get github.com/envoyproxy/protoc-gen-validate && \
    go env -w GO111MODULE=""

RUN go install github.com/envoyproxy/protoc-gen-validate@v0.6.7

COPY gen_all.sh /bin/gen_all.sh

WORKDIR /

ENTRYPOINT ["/bin/gen_all.sh"]
