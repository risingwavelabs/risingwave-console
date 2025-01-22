# Wavekit

Wavekit is a platform for accessing and managing RisingWave clusters. 

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

### Quick start with Docker (for testing)

The wavekit server requires a PostgreSQL database to store the data. The pgbundle version of the wavekit server is a docker image that bundles the PostgreSQL database and the wavekit server. Note that this is NOT recommended for production, as it does not support multiple nodes.

```shell
docker run -p 8020:8020 cloudcarver/wavekit:v0.1.0-pgbundle
```

To persist data, you can use a volume:

```shell
docker run -p 8020:8020 -v wavekit-data:/var/lib/postgresql cloudcarver/wavekit:v0.1.0-pgbundle
```

#### Networking issue with docker setup

Since the wavekit server is running in a docker container, you need to make sure your RisingWave clusters can be accessed inside that wavekit server container.

##### Option 1: Use bridge network

```shell
docker network create wavekit-network
docker run -d --name risingwave --net=test risingwavelabs/risingwave:v2.1.0
docker run -d --name wavekit-server --network wavekit-network -p 8020:8020 cloudcarver/wavekit:v0.1.0-pgbundle
```

With this setup, you can access the RisingWave cluster with the hostname `risingwave`.

##### Option 2: Use the host network (Recommended)

```shell
docker run -d --name risingwave --net=host risingwavelabs/risingwave:v2.1.0
docker run -d --name wavekit-server --net=host -p 8020:8020 cloudcarver/wavekit:v0.1.0-pgbundle
```

With this setup, you can access the RisingWave cluster with the hostname `localhost`.

Note that WSL2 does not support the host network, so this option is not recommended for WSL2.

### Docker compose

TODO

### Kubernetes deployment

TODO


### Binary

TODO

## Configuration


