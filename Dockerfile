# Build the manager binary
FROM golang:1.15.7 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

FROM builder AS controller-builder
COPY api/ api/
COPY controllers/ controllers/
COPY main.go .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -a -o bin/manager main.go

FROM gcr.io/distroless/static:nonroot AS controller
WORKDIR /
COPY --from=controller-builder /workspace/bin/manager .
USER nonroot:nonroot
ENTRYPOINT ["/manager"]

FROM builder AS runner-builder
COPY api/ api/
COPY runner/ runner/
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -a -o bin/runner ./runner

FROM gcr.io/distroless/static:nonroot AS runner
WORKDIR /
COPY runtimes runtimes
COPY --from=runner-builder /workspace/bin/runner .
USER nonroot:nonroot
ENTRYPOINT ["/runner"]
