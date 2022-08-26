FROM golang

RUN apk update && \
    apk add protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

COPY gen_all.sh /bin/gen_all.sh

WORKDIR /

ENTRYPOINT ["/bin/gen_all.sh"]
