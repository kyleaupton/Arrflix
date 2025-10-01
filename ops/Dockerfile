# --- Build API ---
FROM golang:1.22 AS api-build
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
RUN cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

# --- Build Web (Vue SPA) ---
FROM node:20 AS web-build
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# --- Runtime (Nginx + API + Postgres) ---
FROM debian:stable-slim

# install deps: nginx, curl, dumb-init, and postgres
RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx postgresql postgresql-contrib ca-certificates dumb-init curl && \
    rm -rf /var/lib/apt/lists/*

# Copy built artifacts
COPY --from=web-build /web/dist /usr/share/nginx/html
COPY ops/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=api-build /out/api /usr/local/bin/api
COPY ops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Prepare Postgres data directory
RUN mkdir -p /var/lib/postgresql/data && chown -R postgres:postgres /var/lib/postgresql
VOLUME ["/var/lib/postgresql/data", "/data"]

EXPOSE 8484

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]
