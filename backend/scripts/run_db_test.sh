#!/bin/bash

# Set environment variables if needed
# export POSTGRES_HOST="172.16.168.150"
# export POSTGRES_PORT="55432"
# export POSTGRES_DB="postgres"
# export POSTGRES_USER="postgres"
# export POSTGRES_PASSWORD="postgres"

# Navigate to the scripts directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Run the Go script
echo "Running database test script..."
go run db_inspector.go

# Exit with the Go script's exit code
exit $? 