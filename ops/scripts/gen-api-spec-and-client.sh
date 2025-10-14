#/bin/bash
set -euo pipefail

project_root=$(git rev-parse --show-toplevel)

cd $project_root/backend
swag init -g internal/http/http.go -o internal/http/docs --requiredByDefault
# --parseDependencyLevel 1

cd $project_root/web
npm run openapi-ts
