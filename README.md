# Wavekit

Wavekit is a self-hosted platform for accessing and managing RisingWave clusters. 

## Quick start with Docker (for testing)

The WaveKit server requires a PostgreSQL database to store the data. The pgbundle version of the WaveKit server is a docker image that bundles the PostgreSQL database and the WaveKit server. Note that this is NOT recommended for production, as it does not support multiple nodes.

```shell
# 1. start a risingwave instance for testing
docker run --rm -p 4566:4566 -p 5690:5690 -p 5691:5691 risingwavelabs/risingwave:v2.1.2  

# 2. start the WaveKit server in the host network 
sudo docker run -p 8020:8020 --net=host --name wavekit risingwavelabs/wavekit:v0.1.2-pgbundle
```

*Note that rootless docker might not be able to use the host network directly, which make it not able to connect to the RisingWave cluster exposed on the host network.*

To persist data, you can use a volume:

```shell
docker run -p 8020:8020 --net=host --name wavekit -v wavekit-data:/var/lib/postgresql risingwavelabs/wavekit:v0.1.2-pgbundle
``` 

## Deploy

WaveKit uses PostgreSQL as the backend database. You need to deploy a PostgreSQL instance for WaveKit to work.

Here is a sample docker-compose file:

```yaml
version: "3.9"
services:
  wavekit:
    image: cloudcarver/wavekit:v0.1.2
    ports:
      - "8020:8020"
    environment:
      WK_PORT: 8020
      WK_PG_HOST: localhost
      WK_PG_PORT: 5432
      WK_PG_USER: postgres
      WK_PG_PASSWORD: postgres
      WK_PG_DB: postgres
      WK_JWT_SECRET: 9138e41195112b568e22480f18a42dd69b38fab5ee1a36fbf63d49b22097d22a
      WK_ROOT_PASSWORD: 123456
      WK_RISECTLDIR: /

  db: 
    image: "postgres:latest"
    ports:
      - "5432:5432"
    environment:
      POSTGRES_PASSWORD: postgres
    volumes:
      - db-data:/var/lib/postgresql/data

  rw:
    image: risingwavelabs/risingwave:v2.1.2
    ports:
      - 4566:4566
      - 5690:5690
      - 5691:5691
    command: "standalone --meta-opts=" 
                    --listen-addr 0.0.0.0:5690 
                    --advertise-addr rw:5690 
                    --dashboard-host 0.0.0.0:5691 
                    --prometheus-host 0.0.0.0:1250 
                    --backend sqlite  
                    --sql-endpoint /root/single_node.db 
                    --state-store hummock+fs:///root/state_store 
                    --data-directory hummock_001" 
                 --compute-opts=" 
                    --listen-addr 0.0.0.0:5688 
                    --prometheus-listener-addr 0.0.0.0:1250 
                    --advertise-addr rw:5688 
                    --async-stack-trace verbose 
                    --parallelism 4 
                    --total-memory-bytes 2147483648 
                    --role both 
                    --meta-address http://0.0.0.0:5690" 
                 --frontend-opts=" 
                   --listen-addr 0.0.0.0:4566 
                   --advertise-addr rw:4566 
                   --prometheus-listener-addr 0.0.0.0:1250 
                   --health-check-listener-addr 0.0.0.0:6786 
                   --meta-addr http://0.0.0.0:5690 
                   --frontend-total-memory-bytes=500000000" 
                 --compactor-opts=" 
                   --listen-addr 0.0.0.0:6660 
                   --prometheus-listener-addr 0.0.0.0:1250 
                   --advertise-addr rw:6660 
                   --meta-address http://0.0.0.0:5690 
                   --compactor-total-memory-bytes=1000000000""

volumes:
  db-data:
  wavekit-data:

```

## Configuration

The WaveKit server can be configured by a YAML file. Mount the `config.yaml` file to `/app/config.yaml` in the container.

```yaml
init: string
port: integer
jwt:
  secret: string
  randomsecret: true/false
pg:
  host: string
  user: string
  password: string
  db: string
  port: integer
root:
  password: string
nointernet: true/false
risectldir: string

```

The configuration can be overridden by the environment variables:

| Environment Variable | Expected Value | Description |
|---------------------|----------------|-------------|
| `WK_INIT` | `string` | The path of file to store the initialization data |
| `WK_PORT` | `integer` | The port of the wavekit server |
| `WK_JWT_SECRET` | `string` | The secret of the jwt |
| `WK_JWT_RANDOMSECRET` | `true/false` | Whether to use a random secret |
| `WK_PG_HOST` | `string` | The host of the postgres database |
| `WK_PG_USER` | `string` | The user of the postgres database |
| `WK_PG_PASSWORD` | `string` | The password of the postgres database |
| `WK_PG_DB` | `string` | The database of the postgres database |
| `WK_PG_PORT` | `integer` | The port of the postgres database |
| `WK_ROOT_PASSWORD` | `string` | The password of the root user, if not set, the default password is "123456" |
| `WK_NOINTERNET` | `true/false` | Whether to disable internet access |
| `WK_RISECTLDIR` | `string` | The path of the directory to store the risectl files |


## Initialization Data

You can use init file to initialize the WaveKit server.

```yaml
clusters:
  - name: Default Local Cluster
    version: v2.1.2
    connections:
      host: rw
      sqlPort: 4566
      metaPort: 5690
      httpPort: 5691
databases:
  - name: rw
    cluster: Default Local Cluster
    username: root
    database: dev

```

## High Availability

The WaveKit server is a stateless service, you can deploy multiple instances of the WaveKit server and use a load balancer to route the requests to the WaveKit server. 
