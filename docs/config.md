# Configuring WaveKit

## Configuration File

WaveKit is configured using a YAML file. This file includes several sections that define the settings for your WaveKit environment:

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

## Overriding Configuration with Environment Variables

You can override the YAML configuration settings by using environment variables. The following environment variables are supported:

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


# Automated Initialization

WaveKit supports automated initialization through a configuration file, eliminating the need for manual setup through the web UI. This feature enables you to programmatically configure your WaveKit environment by pre-defining cluster and database information.

Here's an example initialization file structure:

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

To use the initialization file, start WaveKit with the `WK_INIT` environment variable pointing to your file:

```shell
WK_INIT=/path/to/init.yaml wavekit
```
