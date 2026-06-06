// rust-app — Docker Tutorial sample (Multi-stage Dockerfile only, no Compose)
//
// Demonstrates:
//   - A two-stage Dockerfile: rust builder → debian-slim runtime
//   - cargo build --release produces an optimised binary
//   - .dockerignore keeps the huge target/ directory out of the build context
//   - No Compose needed — a single `docker build` + `docker run` workflow
//
// Build and run with Docker:
//
//   docker build -t rust-app .
//   docker run --rm -p 3000:3000 rust-app
//   curl http://localhost:3000
//   curl http://localhost:3000/health

use axum::{Json, Router, routing::get};
use serde_json::{Value, json};
use std::net::SocketAddr;

async fn root() -> &'static str {
    "<h1>Hello from Docker (Rust)!</h1>\
     <p>This binary was compiled with <strong>cargo build --release</strong> \
     inside a Rust builder stage and copied into a slim Debian runtime image.</p>\
     <p>No Compose needed — just <code>docker build</code> + <code>docker run</code>.</p>"
}

async fn health() -> Json<Value> {
    Json(json!({ "status": "ok", "runtime": "rust" }))
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/", get(root))
        .route("/health", get(health));

    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    println!("Listening on http://{}", addr);

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
