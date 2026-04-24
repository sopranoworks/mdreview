#!/bin/bash
# mdreview smart wrapper
# Automatically chooses between binary and source execution.

# Get the directory where this script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

if [ -f "./mdreview" ]; then
    # Use pre-built binary for performance
    exec ./mdreview "$@"
else
    # Fallback to running from source for zero-setup
    if command -v go >/dev/null 2>&1; then
        exec go run . "$@"
    else
        echo "Error: mdreview binary not found and 'go' is not installed." >&2
        exit 1
    fi
fi
