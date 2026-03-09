# Stage 1 - build
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/paoloelefante/perf-agent/internal/version.Version=${VERSION}" \
    -o /perf-agent \
    ./cmd/perf-agent

# Stage 2 - runtime
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /perf-agent /perf-agent

EXPOSE 8080

ENTRYPOINT ["/perf-agent"]
