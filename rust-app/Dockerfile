# rust-app/Dockerfile
# Pattern: multi-stage build — large Rust builder → slim Debian runtime
#           Run with plain docker commands (no Compose)
#
# Key teaching points:
#   1. Stage 1 (builder) has the full Rust toolchain + cargo — needed to compile
#   2. cargo build --release produces an optimised binary (~5–15 MB)
#   3. Stage 2 (runtime) is debian:bookworm-slim — glibc available, shell available
#      (unlike distroless) — useful when you need basic debugging tools
#   4. COPY --from=builder brings over only the release binary
#   5. .dockerignore keeps the local target/ directory (can be gigabytes) out of
#      the build context — critical for fast builds

# ── Stage 1: Builder ──────────────────────────────────────────────────────────
FROM rust:1-slim AS builder
WORKDIR /build
# Copy manifests first to cache the dependency compilation layer.
# Only reruns `cargo build` for deps when Cargo.toml / Cargo.lock changes.
COPY Cargo.toml ./
# Create a dummy src/main.rs so cargo can compile dependencies alone
RUN mkdir src && echo 'fn main(){}' > src/main.rs
RUN cargo build --release
# Now overwrite with real source and build the actual binary
COPY src ./src
# Touch main.rs so cargo knows it changed (otherwise it skips it)
RUN touch src/main.rs && cargo build --release

# ── Stage 2: Runtime ─────────────────────────────────────────────────────────
# debian:bookworm-slim provides glibc + basic utilities but strips most packages.
# Contrast with distroless (no shell at all) — choose based on your ops needs.
FROM debian:bookworm-slim
# Install CA certificates in case the app makes HTTPS calls
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
# Create a non-root user
RUN useradd --create-home appuser
USER appuser
COPY --from=builder /build/target/release/rust-app /usr/local/bin/rust-app
EXPOSE 3000
CMD ["rust-app"]

