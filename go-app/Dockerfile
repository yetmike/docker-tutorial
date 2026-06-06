# go-app/Dockerfile
# Pattern: multi-stage build — large builder → tiny distroless runtime
#
# Key teaching points:
#   1. Stage 1 (builder) has the full Go toolchain — needed to compile, never shipped
#   2. CGO_ENABLED=0 produces a fully static binary (no glibc dependency)
#   3. Stage 2 (runtime) is distroless/static — no shell, no OS packages, ~2 MB base
#   4. COPY --from=builder copies ONLY the compiled binary
#   5. Result: a production image that is tiny and has minimal attack surface

# ── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM golang:1.22 AS builder

WORKDIR /src

# Download dependencies first (cached layer — only reruns if go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a fully static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app .

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
# distroless/static is a Google-maintained image with NO shell, NO package manager.
# It contains only CA certificates, timezone data, and glibc-less runtime support.
FROM gcr.io/distroless/static:nonroot

# Copy the compiled binary from the builder stage
COPY --from=builder /app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
