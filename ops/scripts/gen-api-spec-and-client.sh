#/bin/bash
set -euo pipefail

project_root=$(pwd)

cd $project_root/backend/internal/http
swag init

cd $project_root/web
npm run openapi-ts
