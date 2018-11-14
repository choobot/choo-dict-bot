FROM golang:1.11.2-alpine3.8
ENV SRC_DIR /go/src/app/
WORKDIR ${SRC_DIR}
COPY app/ ${SRC_DIR}
RUN apk add build-base && \
    apk add git && \
    cd ${SRC_DIR} && \
    go get ./... && \
    go install -v
ENTRYPOINT [ "/bin/sh","-c" ]
CMD [ "/go/bin/app" ]