#!/bin/bash
set -e

SETUP_SCRIPT="ops/setup_streams.sh"

echo "Verifying Stream Setup script..."

if [ ! -f "$SETUP_SCRIPT" ]; then
    echo "❌ Error: $SETUP_SCRIPT does not exist."
    exit 1
fi

if ! grep -q "nats stream add EVENTS" "$SETUP_SCRIPT"; then
    echo "❌ Error: Stream creation command not found in $SETUP_SCRIPT."
    exit 1
fi

if ! grep -q "enterprise.>" "$SETUP_SCRIPT"; then
    echo "❌ Error: Correct subject not found in $SETUP_SCRIPT."
    exit 1
fi

echo "✅ Stream Setup script verification passed."
exit 0
