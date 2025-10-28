# Snaggle Ops - Container Orchestration

Snaggle Ops is a Go-based reconciliation controller that manages Docker containers for the Snaggle media management platform. It replaces the previous supervisord-based approach with a Kubernetes-style controller pattern.

## Architecture

The reconciler operates as a single container (`snaggle-ops`) that spawns and manages sibling Docker containers via the Docker socket. It handles both static services (always running) and dynamic services (user-configurable).

### Static Services

These services are always defined and running:

- **PostgreSQL** (`snaggle-postgres`) - Database
- **Prowlarr** (`snaggle-prowlarr`) - Indexer management
- **API** (`snaggle-api`) - Snaggle backend API
- **Nginx** (`snaggle-nginx`) - Web server and reverse proxy

### Dynamic Services

These services are created based on database configuration:

- **qBittorrent** - Torrent client instances
- **Transmission** - Alternative torrent client (planned)

## Directory Structure

```
ops/
├── go.mod                           # Go module for snaggle-ops
├── go.sum
├── main.go                          # Entry point
├── Dockerfile                       # Builds snaggle-ops binary
├── docker-compose.yml               # Single service: snaggle-ops
├── internal/
│   ├── config/
│   │   └── config.go               # Config from env vars
│   ├── reconciler/
│   │   └── reconciler.go           # Core reconcile loop
│   ├── docker/
│   │   └── client.go               # Docker API wrapper
│   └── services/
│       ├── service.go              # Service interface
│       ├── postgres.go             # Static service
│       ├── nginx.go                # Static service
│       ├── prowlarr.go             # Static service
│       ├── api.go                  # Static service
│       ├── dynamic.go              # Dynamic service factory
│       ├── qbittorrent.go          # Dynamic service impl
│       └── transmission.go         # Dynamic service impl
└── images/
    ├── Dockerfile.api              # API service image
    ├── Dockerfile.nginx            # Nginx service image
    └── Dockerfile.prowlarr         # Prowlarr service image
```

## Usage

### Running with Docker Compose

```bash
cd ops
docker-compose up -d
```

This will start the `snaggle-ops` container which will then spawn and manage all other services.

### Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `POSTGRES_DB` - Database name (default: snaggle)
- `POSTGRES_USER` - Database user (default: snaggle)
- `POSTGRES_PASSWORD` - Database password (default: snaggle)
- `JWT_SECRET` - JWT signing secret
- `TMDB_API_KEY` - TMDB API key
- `RECONCILE_INTERVAL` - Reconciliation interval (default: 10s)
- `NETWORK_NAME` - Docker network name (default: snaggle-network)

## Database Schema

The reconciler uses a `service_instance` table to manage dynamic services:

```sql
CREATE TABLE service_instance (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL, -- 'qbittorrent', 'transmission', etc.
  enabled BOOLEAN NOT NULL DEFAULT true,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Adding Dynamic Services

To add a qBittorrent instance:

```sql
INSERT INTO service_instance (name, type, config) VALUES (
  'my-qbittorrent',
  'qbittorrent',
  '{
    "webui_port": "8080",
    "config_path": "/host/path/to/config",
    "downloads_path": "/host/path/to/downloads"
  }'::jsonb
);
```

## Reconciliation Loop

The reconciler runs in a continuous loop:

1. **Ensure Infrastructure**: Create network and volumes if needed
2. **Build Desired State**: Combine static services with enabled dynamic services from DB
3. **Get Actual State**: Query Docker for running managed containers
4. **Reconcile**: Create, update, or remove containers to match desired state
5. **Handle Dependencies**: Ensure services start in correct order
6. **Health Check**: Wait for services to be healthy before starting dependents

## Service Dependencies

Services have explicit dependencies:

- `snaggle-api` depends on `snaggle-postgres` and `snaggle-prowlarr`
- `snaggle-nginx` depends on `snaggle-api`
- Dynamic services (qBittorrent, etc.) have no dependencies

## Container Labels

All managed containers receive these labels:

- `snaggle.managed=true` - Identifies Snaggle-managed containers
- `snaggle.service=<name>` - Service identifier
- `snaggle.type=<type>` - Service type (postgres, api, qbittorrent, etc.)
- `snaggle.instance=<name>` - Instance name (for dynamic services)

## Building Service Images

To build the individual service images:

```bash
# Build API image
docker build -f images/Dockerfile.api -t snaggle-api:latest ..

# Build Nginx image
docker build -f images/Dockerfile.nginx -t snaggle-nginx:latest ..

# Build Prowlarr image
docker build -f images/Dockerfile.prowlarr -t snaggle-prowlarr:latest .
```

## Migration from Supervisord

The previous supervisord-based approach has been moved to `ops-old/` for reference. Key changes:

- Single container approach → Multi-container orchestration
- Process management → Container lifecycle management
- Shell scripts → Go reconciliation logic
- Static configuration → Database-driven dynamic services

## Benefits

- **Single Entry Point**: Users only need `snaggle-ops` in docker-compose
- **Dynamic Scaling**: Add/remove services via database without restarts
- **Better Isolation**: Each service in its own container
- **Easier Debugging**: Individual container logs and inspection
- **Testable**: Reconciler logic can be unit tested
- **Kubernetes-Ready**: Similar pattern to k8s controllers

## Development

### Building the Reconciler

```bash
cd ops
go build -o snaggle-ops .
```

### Running Locally

```bash
# Ensure Docker socket is accessible
sudo chmod 666 /var/run/docker.sock

# Run the reconciler
./snaggle-ops
```

### Testing

The reconciler can be tested by:

1. Starting with `docker-compose up`
2. Adding dynamic services via database
3. Verifying containers are created and managed
4. Checking service dependencies and health

## Troubleshooting

### Container Not Starting

Check the reconciler logs:

```bash
docker logs snaggle-ops
```

### Service Dependencies

Verify dependencies are healthy:

```bash
docker ps --filter "label=snaggle.managed=true"
```

### Database Connection

Ensure the database URL is correct and PostgreSQL is accessible from the reconciler container.

## Future Enhancements

- Config hash comparison for efficient updates
- Proper topological sorting for dependencies
- Metrics and monitoring endpoints
- Service discovery and health checking
- Rolling updates for service changes
