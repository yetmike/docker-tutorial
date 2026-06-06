# flask-app/Dockerfile
# Pattern: single-stage Dockerfile for an interpreted runtime
#
# Key teaching points:
#   1. Official slim base image keeps the image small without a multi-stage build
#   2. Copy requirements BEFORE source code so the pip install layer is cached
#      between code-only changes (dependencies change far less often)
#   3. EXPOSE documents the port; it doesn't actually publish it (compose.yaml does)

FROM python:3.12-slim

# Create a non-root user for security best practice
RUN adduser --disabled-password --gecos "" appuser

WORKDIR /app

# Copy and install dependencies first — this layer is cached until
# requirements.txt changes, even if app.py changes.
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application source
COPY app.py .

# Switch to non-root user
USER appuser

EXPOSE 5000

CMD ["python", "app.py"]
