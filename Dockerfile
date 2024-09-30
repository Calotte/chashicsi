FROM golang:1.21.13

ARG GOPROXY

WORKDIR /workspace
COPY . .

RUN make build
RUN chmod u+x /workspace/bin/manager

ENTRYPOINT ["/workspace/bin/manager"]
