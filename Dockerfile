FROM golang:1.13 AS builder
RUN mkdir -p /src/github.com/ibrokethecloud/rio-sample
COPY . /src/github.com/ibrokethecloud/rio-sample
RUN cd /src/github.com/ibrokethecloud/rio-sample \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /root/rio-sample

## Using upstream aquasec kube-bench and layering it up
FROM scratch
ARG COLOR
ENV COLOR=$COLOR
COPY --from=builder /root/rio-sample /rio-sample
WORKDIR /

ENTRYPOINT ["/rio-sample"]
