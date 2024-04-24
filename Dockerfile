# Build the manager binary
FROM registry.cn-hangzhou.aliyuncs.com/godev/golang:1.22 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
ENV GOPROXY https://goproxy.cn,direct
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/controller/ internal/controller/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.cn-hangzhou.aliyuncs.com/godev/alpine:latest
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENV MCENTER_GRPC_ADDRESS=127.0.0.1:18010 \
MCENTER_CLINET_ID=moperator \
MCENTER_CLIENT_SECRET=moperator

ENTRYPOINT ["/manager"]
