# --- Build API ---
FROM golang:1.24 AS api-build
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

# --- Runtime (Nginx + API + Postgres + Supervisor) ---
FROM debian:stable-slim AS prod-runtime

# install deps: nginx, curl, dumb-init, postgres, and supervisor
RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx postgresql postgresql-contrib ca-certificates dumb-init curl supervisor wget xz-utils tar jq && \
    rm -rf /var/lib/apt/lists/*

# Make Postgres client/server binaries easy to find on PATH
# Debian installs them under /usr/lib/postgresql/<major>/bin
RUN set -eu; \
    pg_bin_base="/usr/lib/postgresql"; \
    if [ -d "$pg_bin_base" ]; then \
      for vdir in "$pg_bin_base"/*; do \
        if [ -d "$vdir/bin" ]; then \
          ln -sf "$vdir/bin/initdb" /usr/local/bin/initdb || true; \
          ln -sf "$vdir/bin/pg_ctl" /usr/local/bin/pg_ctl || true; \
          ln -sf "$vdir/bin/postgres" /usr/local/bin/postgres || true; \
          ln -sf "$vdir/bin/psql" /usr/local/bin/psql || true; \
        fi; \
      done; \
    fi

# Copy built artifacts
COPY --from=web-build /web/dist /usr/share/nginx/html
COPY ops/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=api-build /out/api /usr/local/bin/api
COPY ops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Supervisor + helper scripts
COPY ops/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ops/scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh
COPY ops/scripts/run-postgres.sh /usr/local/bin/run-postgres.sh
RUN chmod +x /usr/local/bin/wait-for-db.sh /usr/local/bin/run-postgres.sh
COPY ops/scripts/wait-for-jackett.sh /usr/local/bin/wait-for-jackett.sh
RUN chmod +x /usr/local/bin/wait-for-jackett.sh
COPY ops/scripts/jacket-print-apikey.sh /usr/local/bin/jacket-print-apikey.sh
RUN chmod +x /usr/local/bin/jacket-print-apikey.sh

# Install Jackett (Linux AMD64). Place under /opt/jackett
RUN set -eux; \
    mkdir -p /opt/jackett && cd /opt/jackett && \
    curl -fsSL https://github.com/Jackett/Jackett/releases/download/v0.24.131/Jackett.Binaries.LinuxAMDx64.tar.gz -o jackett.tar.gz && \
    tar -xzf jackett.tar.gz --strip-components=1 && \
    rm jackett.tar.gz && \
    ln -sf /opt/jackett/jackett /usr/local/bin/jackett

# Prepare Postgres data directory
RUN mkdir -p /var/lib/postgresql/data && chown -R postgres:postgres /var/lib/postgresql
VOLUME ["/var/lib/postgresql/data", "/data"]

EXPOSE 8484

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]

# --- Dev Runtime (Nginx + API live reload + Postgres + Supervisor + Vite) ---
FROM golang:1.24 AS dev-runtime

# install deps: nginx, curl, dumb-init, postgres, supervisor, build tools, node
RUN apt-get update && \
    apt-get install -y --no-install-recommends nginx postgresql postgresql-contrib ca-certificates dumb-init curl supervisor git build-essential wget xz-utils tar jq && \
    rm -rf /var/lib/apt/lists/*

# Node 20 via NodeSource
RUN set -eux; \
    arch=$(dpkg --print-architecture); \
    case "$arch" in \
      amd64) node_arch=x64 ;; \
      arm64) node_arch=arm64 ;; \
      *) echo "Unsupported arch: $arch"; exit 1 ;; \
    esac; \
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get update && apt-get install -y nodejs && \
    rm -rf /var/lib/apt/lists/*

# Make Postgres binaries available on PATH
RUN set -eu; \
    pg_bin_base="/usr/lib/postgresql"; \
    if [ -d "$pg_bin_base" ]; then \
      for vdir in "$pg_bin_base"/*; do \
        if [ -d "$vdir/bin" ]; then \
          ln -sf "$vdir/bin/initdb" /usr/local/bin/initdb || true; \
          ln -sf "$vdir/bin/pg_ctl" /usr/local/bin/pg_ctl || true; \
          ln -sf "$vdir/bin/postgres" /usr/local/bin/postgres || true; \
          ln -sf "$vdir/bin/psql" /usr/local/bin/psql || true; \
        fi; \
      done; \
    fi

# Install Air for Go live reload
RUN curl -fsSL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | bash -s -- -b /usr/local/bin

WORKDIR /app

# Copy configs and scripts
COPY ops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
COPY ops/nginx/default.dev.conf /etc/nginx/conf.d/default.dev.conf
COPY ops/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY ops/supervisor/supervisord.dev.conf /etc/supervisor/conf.d/supervisord.dev.conf
COPY ops/supervisor/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ops/scripts/wait-for-db.sh /usr/local/bin/wait-for-db.sh
COPY ops/scripts/run-postgres.sh /usr/local/bin/run-postgres.sh
RUN chmod +x /usr/local/bin/wait-for-db.sh /usr/local/bin/run-postgres.sh
COPY ops/scripts/wait-for-jackett.sh /usr/local/bin/wait-for-jackett.sh
RUN chmod +x /usr/local/bin/wait-for-jackett.sh
COPY ops/scripts/jacket-print-apikey.sh /usr/local/bin/jacket-print-apikey.sh
RUN chmod +x /usr/local/bin/jacket-print-apikey.sh

# Install Jackett (Linux AMD64) into dev image as well
RUN set -eux; \
    mkdir -p /opt/jackett && cd /opt/jackett && \
    curl -fsSL https://github.com/Jackett/Jackett/releases/download/v0.24.131/Jackett.Binaries.LinuxAMDx64.tar.gz -o jackett.tar.gz && \
    tar -xzf jackett.tar.gz --strip-components=1 && \
    rm jackett.tar.gz && \
    ln -sf /opt/jackett/jackett /usr/local/bin/jackett

# Prepare dirs and volumes
RUN mkdir -p /var/lib/postgresql/data && chown -R postgres:postgres /var/lib/postgresql
VOLUME ["/var/lib/postgresql/data", "/data"]

# Placeholders for source mounts; we still need them for context
WORKDIR /app

EXPOSE 8484
EXPOSE 9117

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/entrypoint.sh"]
