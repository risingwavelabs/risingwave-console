# Wavekit

Wavekit is a self-hosted platform for accessing and managing RisingWave clusters. 

## Quick start with Docker (for testing)

The WaveKit server requires a PostgreSQL database to store the data. The pgbundle version of the WaveKit server is a docker image that bundles the PostgreSQL database and the WaveKit server. Note that this is NOT recommended for production, as it does not support multiple nodes.

```shell
# 1. start a risingwave instance for testing
docker run --rm -p 4566:4566 -p 5690:5690 -p 5691:5691 risingwavelabs/risingwave:v2.1.2  

# 2. start the WaveKit server in the host network 
sudo docker run --net=host --name wavekit risingwavelabs/wavekit:v0.1.2-pgbundle

# 3. Open the browser and go to http://localhost:8020, the default username and password are "root" and "root".
```

*Note that rootless docker might not be able to use the host network directly, which make it not able to connect to the RisingWave cluster exposed on the host network.*

To persist data, you can use a volume:

```shell
docker run --net=host --name wavekit -v wavekit-data:/var/lib/postgresql risingwavelabs/wavekit:v0.1.2-pgbundle
``` 

Using the pgbundle version in production is NOT recommended, as it integrates the Postgres process with default configurations into the image.

## Deploy

WaveKit uses PostgreSQL as the backend database. You need to deploy a PostgreSQL instance for WaveKit to work.

Here is a sample docker-compose file:

```yaml
{{README_DOCKER_COMPOSE}}
```

## Configuration

The WaveKit server can be configured by a YAML file. Mount the `config.yaml` file to `/app/config.yaml` in the container.

```yaml
{{README_YAML}}
```

The configuration can be overridden by the environment variables:

{{README_ENV}}

## Initialization Data

You can use init file to initialize the WaveKit server.

```yaml
{{README_INIT}}
```

## High Availability

The WaveKit server is a stateless service, you can deploy multiple instances of the WaveKit server and use a load balancer to route the requests to the WaveKit server. 
