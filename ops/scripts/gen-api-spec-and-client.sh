#/bin/bash
set -euo pipefail

# better way of atomically getting the project root?
project_root=$(git rev-parse --show-toplevel)

# ensure we're in the backend directory
cd $project_root/backend

swag init -g internal/http/http.go -o internal/http/docs

cd $project_root/web
npm run openapi-ts
