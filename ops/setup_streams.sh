#!/bin/sh
set -e

echo "Creating EVENTS stream..."
nats stream add EVENTS --subjects "enterprise.>" --storage file --retention limits --max-msgs=-1 --max-bytes=-1 --max-age=1y --discard old --dupe-window=2m --replicas=1 --no-ack

echo "Stream EVENTS created."
