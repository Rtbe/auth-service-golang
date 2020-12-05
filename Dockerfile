# This Dockerfile specifies heroku deployment.
FROM golang:1.15-buster AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN go build -v -o /bin/server cmd/*.go

FROM ubuntu:20.10 as base

# MongoDB download URL
ARG DB_URL=https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu1804-4.2.6.tgz

# Server ENVs.
ENV MODE=heroku
ENV TOKEN_SECRET=herokusecret 
ENV DB_USER=admin
ENV DB_PASSWORD=password
ENV DB_NAME=testTask
ENV DB_PORT=27017

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl && \
    curl -OL ${DB_URL} && \
    tar -zxvf mongodb-linux-x86_64-ubuntu1804-4.2.6.tgz && \
    mv ./mongodb-linux-x86_64-ubuntu1804-4.2.6/bin/* /usr/local/bin/ && \
    rm -rf ./mongodb-linux-x86_64-ubuntu1804-4.2.6 && rm ./mongodb-linux-x86_64-ubuntu1804-4.2.6.tgz

COPY ./scripts/init-mongodbs.sh ./scripts/init-replica.sh ./scripts/entry-point.sh /

RUN chmod +x /init-mongodbs.sh && \
    chmod +x /init-replica.sh && \
    chmod +x /entry-point.sh

# Data directory
ARG DB1_DATA_DIR=/var/lib/mongo1
ARG DB2_DATA_DIR=/var/lib/mongo2
ARG DB3_DATA_DIR=/var/lib/mongo3

# Log directory
ARG DB1_LOG_DIR=/var/log/mongodb1
ARG DB2_LOG_DIR=/var/log/mongodb2
ARG DB3_LOG_DIR=/var/log/mongodb3

# DB Ports
ARG DB1_PORT=27017
ARG DB1_PORT=27018
ARG DB1_PORT=27019

RUN mkdir -p ${DB1_DATA_DIR} && \
    mkdir -p ${DB1_LOG_DIR} && \
    mkdir -p ${DB2_DATA_DIR} && \
    mkdir -p ${DB2_LOG_DIR} && \
    mkdir -p ${DB3_DATA_DIR} && \
    mkdir -p ${DB3_LOG_DIR} && \
    chown `whoami` ${DB1_DATA_DIR} && \
    chown `whoami` ${DB1_LOG_DIR} && \
    chown `whoami` ${DB2_DATA_DIR} && \
    chown `whoami` ${DB2_LOG_DIR} && \
    chown `whoami` ${DB3_DATA_DIR} && \
    chown `whoami` ${DB3_LOG_DIR}

EXPOSE ${DB1_PORT}
EXPOSE ${DB2_PORT}
EXPOSE ${DB3_PORT}

COPY --from=build /bin/server /

ENTRYPOINT [ "bash", "entry-point.sh" ]