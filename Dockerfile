# Build the manager binary
FROM golang:1.12.5 as builder

WORKDIR /workspace

ENV TELEPORT_VERSION=4.0.2
RUN curl -sSLf https://get.gravitational.com/teleport-v${TELEPORT_VERSION}-linux-amd64-bin.tar.gz | tar -xzvf - --strip-component=1 teleport/tctl

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM quay.io/cybozu/ubuntu:18.04
WORKDIR /
COPY --from=builder /workspace/tctl .
COPY --from=builder /workspace/manager .
ENTRYPOINT ["/manager"]
