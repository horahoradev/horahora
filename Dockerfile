FROM golang

RUN apt-get update && \
    apt install -y protobuf-compiler && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

COPY gen_all.sh /bin/gen_all.sh

#COPY video_service/ /video_service/
#COPY scheduler/ /scheduler/
#COPY user_service/ /user_service/


WORKDIR /

ENTRYPOINT ["/bin/gen_all.sh"]