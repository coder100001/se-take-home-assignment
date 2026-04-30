#!/bin/bash

echo "Running CLI application..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

if [ ! -f order-controller ]; then
    echo "Building first..."
    go build -o order-controller ./cmd/main.go
fi

printf 'normal burger\nvip fries\n+bot\n+bot\nstatus\nquit\n' | ./order-controller > "$SCRIPT_DIR/result.txt"

echo "CLI application execution completed"
