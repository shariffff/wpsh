# Build stage
FROM golang:1.24-alpine AS builder

# Install git for version info
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG COMMIT_SHA=unknown
ARG BUILD_DATE=unknown
ARG GO_OS=linux
ARG GO_ARCH=amd64

# Build the binary
RUN CGO_ENABLED=0 GOOS=${GO_OS} GOARCH=${GO_ARCH} go build \
    -ldflags "-X github.com/wordmon/cli/cmd.Version=${VERSION} \
              -X github.com/wordmon/cli/cmd.CommitSHA=${COMMIT_SHA} \
              -X github.com/wordmon/cli/cmd.BuildDate=${BUILD_DATE}" \
    -o wordmon .

# Output stage - just the binary
FROM scratch AS binary
COPY --from=builder /build/wordmon /wordmon
