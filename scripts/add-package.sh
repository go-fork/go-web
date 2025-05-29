#!/bin/bash

# Script to add a new package to the registry

if [ $# -ne 4 ]; then
    echo "Usage: $0 <name> <path> <repository> <description>"
    echo "Example: $0 auth go.fork.vn/auth github.com/go-fork/auth 'Authentication package'"
    exit 1
fi

NAME="$1"
PATH="$2"
REPOSITORY="$3"
DESCRIPTION="$4"

# Create JSON payload
JSON_PAYLOAD=$(cat <<EOF
{
    "name": "$NAME",
    "path": "$PATH",
    "repository": "$REPOSITORY",
    "description": "$DESCRIPTION"
}
EOF
)

# Send POST request to API
curl -X POST \
     -H "Content-Type: application/json" \
     -d "$JSON_PAYLOAD" \
     http://localhost:8080/api/packages

echo
echo "Package $NAME added successfully!"
