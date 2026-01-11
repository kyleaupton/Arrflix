# Arrflix

Arrflix is a self-hosted media management platform designed to unify the best parts of Sonarr, Radarr, and Overseerr into a single, transparent, and flexible tool. It helps you manage your movie and series collections with a focus on filesystem integrity and efficient storage usage.

## Features

- **Unified Management**: Manage both Movies and Series in one interface.
- **Granular Monitoring**: Monitor at the Series, Season, or individual Episode level.
- **Filesystem First**: Prioritizes your filesystem as the source of truth.
- **Efficient Storage**: Uses a hardlink-first strategy for importing downloads to save space.
- **Indexer Integration**: Bundled Prowlarr for seamless indexer management out of the box.
- **External Downloader Support**: Integrates with popular downloaders like qBittorrent.
- **Modular & Secure**: Bring your own downloader and VPN setup for maximum flexibility and privacy.

## Getting Started

Arrflix is distributed as a Docker container that includes the API, Frontend, and Prowlarr.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- An external downloader (e.g., qBittorrent) and your preferred VPN solution.

### Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/kyleaupton/Arrflix.git
   cd Arrflix
   ```

2. **Configure environment variables**:
   Create a `.env` file in the root directory:

   ```env
   TMDB_API_KEY=your_tmdb_api_key_here
   MEDIA_LIBRARIES=/path/to/your/media
   ```

3. **Launch Arrflix**:

   ```bash
   docker compose up -d arrflix
   ```

4. **Access the interface**:
   Open your browser and navigate to `http://localhost:8484`.

## Downloader & VPN Scope

Arrflix focuses on media management and indexer orchestration. To maintain a lightweight and flexible footprint:

- **Downloaders**: Arrflix integrates with external downloaders via API. While sample services are provided in the `docker-compose.yml` for reference, you are encouraged to use your existing downloader setup.
- **VPN/Privacy**: Privacy is your responsibility. We recommend running your downloaders behind a VPN (such as Gluetun or a dedicated VPN container). Arrflix does not bundle a VPN client.

## Configuration

For detailed configuration options, including adding indexers and setting up downloaders, please refer to the settings section within the web interface.

## Documentation

- **User Guide**: Coming soon(ish).
- **Developer Documentation**: See [docs/PROJECT_BRIEF.md](docs/PROJECT_BRIEF.md) for architecture and development setup.

## License

TODO: figure this out
