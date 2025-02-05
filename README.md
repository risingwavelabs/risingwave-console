# Wavekit

Wavekit is a self-hosted platform for accessing and managing RisingWave clusters. 

## Features

v1.0:
- [x] SQL Editor
- [x] Materialized View Viewer
- [x] Web UI for Risectl
- [x] Meta Backup
- [x] Collect Diagnostic Data
- [ ] Auto Diagnostic
- [ ] Auto Backup

v2.0:
- [ ] Monitoring
- [ ] UDF Management
- [ ] SSO
- [ ] Access Control
- [ ] SQL Review

## Get Started

### Quck start with Docker Compose (for testing)

```shell
git clone git@github.com:risingwavelabs/wavekit.git
cd wavekit/examples/docker-compose-test
docker compose up
```

Open http://localhost:8020 in your browser to access the wavekit web UI. 
The default user is `root` and the default password is `123456`.

### Quick start with Docker (for testing)

The wavekit server requires a PostgreSQL database to store the data. The pgbundle version of the wavekit server is a docker image that bundles the PostgreSQL database and the wavekit server. Note that this is NOT recommended for production, as it does not support multiple nodes.

```shell
docker run --net=host cloudcarver/wavekit:v0.1.1-pgbundle
```

To persist data, you can use a volume:

```shell
docker run --net=host -v wavekit-data:/var/lib/postgresql cloudcarver/wavekit:v0.1.1-pgbundle
```

Note that WSL2 and rootless docker may have the host network issue, in that case, you can use the bridge network to connect the wavekit server and the RisingWave cluster.

### Docker compose

A sample docker compose file is provided in the `examples/docker-compose` directory.

```shell
cd wavekit/examples/docker-compose
docker compose up
```

### High Availability

The wavekit server is a stateless service, you can deploy multiple instances of the wavekit server and use a load balancer to route the requests to the wavekit server. 
