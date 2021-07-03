FROM golang:1.15.7-buster

COPY . /backup_service

WORKDIR /backup_service

RUN apt-get update && \
    apt-get install -y lsb-release && \
    sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    apt-get update && \
    apt-get -y install postgresql && \
    go mod vendor && \
    go build --mod=vendor -gcflags "all=-N -l" -o backup_service .

# RUN apt-get update && \
#     apt-get install postgresql-12 && \
#     go build --mod=vendor -gcflags "all=-N -l" -o backup_service .

ENTRYPOINT ["./backup_service"]