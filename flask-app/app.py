"""
flask-app — Docker Tutorial sample (Dockerfile + Compose pattern)

Demonstrates:
  - A single-stage Python Dockerfile
  - Docker Compose wiring two services (Flask + Redis) together
  - Service-name resolution: this app reaches Redis simply as "redis"

Run with Docker Compose:
  docker compose up --build
  curl http://localhost:5000
  curl http://localhost:5000/health
"""

import os
import redis
from flask import Flask, jsonify

app = Flask(__name__)

REDIS_URL = os.environ.get("REDIS_URL", "redis://localhost:6379")
r = redis.from_url(REDIS_URL, decode_responses=True)


@app.route("/")
def index():
    visits = r.incr("visits")
    return (
        f"<h1>Hello from Docker!</h1>"
        f"<p>This page has been visited <strong>{visits}</strong> time(s).</p>"
        f"<p>Visit count is stored in <em>Redis</em> — "
        f"running as a separate container linked via Docker Compose.</p>"
    )


@app.route("/health")
def health():
    try:
        r.ping()
        redis_status = "ok"
    except Exception as exc:
        redis_status = f"error: {exc}"
    return jsonify({"status": "ok", "redis": redis_status})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=False)
