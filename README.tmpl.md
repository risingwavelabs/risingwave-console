# Wavekit

Wavekit is a self-hosted platform for managing RisingWave clusters.

## Quick Start with Docker (for Testing)

The WaveKit server uses a PostgreSQL database to store essential cluster information, including connection details such as hostnames and ports of RisingWave clusters. 

For testing purposes, we provide a `pgbundle` version - a Docker image that conveniently packages both PostgreSQL and the WaveKit server together. **Note:** While this bundled version simplifies testing, it is NOT suitable for production environments since it lacks support for multiple nodes and uses default PostgreSQL configurations.

```shell
# 1. Start a RisingWave instance for testing
docker run --rm -p 4566:4566 -p 5690:5690 -p 5691:5691 risingwavelabs/risingwave:v2.1.2  

# 2. Start the WaveKit server on the host network 
sudo docker run --net=host --name wavekit risingwavelabs/wavekit:v0.1.2-pgbundle

# 3. Open your browser and navigate to http://localhost:8020. 
#    The default username and password are "root" and "root".
```

*Note: Rootless Docker may not be able to use the host network directly, which can prevent it from connecting to the RisingWave cluster exposed on the host network.*

To persist data, use a volume:

```shell
docker run --net=host --name wavekit -v wavekit-data:/var/lib/postgresql risingwavelabs/wavekit:v0.1.2-pgbundle
``` 

Using the `pgbundle` version in production is NOT recommended, as it integrates the Postgres process with default configurations into the image.

## Production Deployment

For production environments, we strongly recommend deploying WaveKit with a dedicated PostgreSQL database for better reliability, scalability and maintainability. The following sections will guide you through setting up WaveKit in a production environment.

WaveKit uses PostgreSQL as its backend database. You need to deploy a PostgreSQL instance for WaveKit to function.

Here is a sample `docker-compose` file:

```yaml
{{README_DOCKER_COMPOSE}}
```

## Configuration

The WaveKit server can be configured using a YAML file. Mount the `config.yaml` file to `/app/config.yaml` in the container.

```yaml
{{README_YAML}}
```

Configuration can be overridden by environment variables:

{{README_ENV}}

## Initialization Data

You can use an init file to initialize the WaveKit server.

```yaml
{{README_INIT}}
```

## High Availability

The WaveKit server is stateless, allowing you to deploy multiple instances and use a load balancer to route requests to the WaveKit server. 
