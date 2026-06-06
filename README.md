# Docker Tutorial

Three runnable sample apps (Python, Go, Rust) — each demonstrating a different
level of Docker packaging.

Every subproject has its own `mise.toml` pinning the exact toolchain.
Run `mise install` inside any subproject to provision its runtime.

---

## Layout

```
docker-tutorial/
├── flask-app/      Python Flask + Redis                     (mise: python 3.12)
├── go-app/         Go + Postgres                            (mise: go 1.22)
└── rust-app/       Rust HTTP server                         (mise: rust stable)
```

---

## The three Docker patterns

| Project | Pattern | What it teaches |
|---------|---------|----------------|
| `flask-app` | **Dockerfile + Compose** | Single-stage image for an interpreted runtime; Compose wires Flask to Redis |
| `go-app` | **Multi-stage + Compose** | Builder → distroless tiny image; Compose adds Postgres |
| `rust-app` | **Multi-stage, no Compose** | Rust builder → debian-slim; `.dockerignore` keeps `target/` out; plain `docker build`/`run` |

---

## Prerequisites

- [mise](https://mise.jdx.dev) (already installed)
- [Docker](https://docs.docker.com/get-docker/) with the Compose plugin

---

## flask-app — Dockerfile + Compose

```bash
cd flask-app
docker compose up --build
# → http://localhost:5000       (visit counter backed by Redis)
# → http://localhost:5000/health

docker compose down
```

Run locally (no Docker):

```bash
mise install            # provisions python 3.12
pip install -r requirements.txt
python app.py
```

---

## go-app — Multi-stage + Compose

```bash
cd go-app
docker compose up --build
# → http://localhost:8080       (Postgres connectivity check)
# → http://localhost:8080/health

# Check how small the final image is (no Go toolchain inside):
docker images go-app-app

docker compose down
```

Run locally (no Docker):

```bash
cd go-app
mise install            # provisions go 1.22
go mod download
go run .
```

---

## rust-app — Multi-stage Dockerfile only

```bash
cd rust-app
docker build -t rust-app .
docker run --rm -p 3000:3000 rust-app
# → http://localhost:3000
# → http://localhost:3000/health
```

Run locally (no Docker):

```bash
mise install            # provisions rust stable
cargo run
```

---

## Further reading

- [Docker docs](https://docs.docker.com)
- [Compose file reference](https://docs.docker.com/compose/compose-file/)
- [Dockerfile best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [distroless images](https://github.com/GoogleContainerTools/distroless)
- [dive — explore image layers](https://github.com/wagoodman/dive)
