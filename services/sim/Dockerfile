FROM python:3.11-slim as builder

WORKDIR /app

# Install system dependencies (if any needed for numpy/gekko compilation)
# apt-get install -y ...

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Final stage
FROM python:3.11-slim

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /usr/local/bin /usr/local/bin

COPY . .

# Environment variables
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app/src:/app/proto

# Run
CMD ["python", "src/main.py"]
