# Configuring WaveKit

## Configuration File

WaveKit is configured using a YAML file. This file includes several sections that define the settings for your WaveKit environment:

```yaml
init: string
host: string
port: integer
jwt:
  secret: string
pg:
  dsn: string
root:
  password: string
nointernet: true/false
risectldir: string
prometheus:
  default_labels: string

```

## Overriding Configuration with Environment Variables

You can override the YAML configuration settings by using environment variables. The following environment variables are supported:

| Environment Variable | Expected Value | Description |
|---------------------|----------------|-------------|
| `WK_INIT` | `string` | (Optional) The path of file to store the initialization data, if not set, skip the initialization |
| `WK_HOST` | `string` | (Optional) The host of the wavekit server, it is used in the API endpoint of the web UI. If not set, the host will be localhost. |
| `WK_PORT` | `integer` | (Optional) The port of the wavekit server, default is 8020 |
| `WK_JWT_SECRET` | `string` | (Optional) The secret of the jwt. If not set, a random secret will be used. |
| `WK_PG_DSN` | `string` | (Required) The DSN (Data Source Name) for postgres database connection. If specified, Host, Port, User, Password, and Db settings will be ignored. |
| `WK_ROOT_PASSWORD` | `string` | (Optional) The password of the root user, if not set, the default password is "123456" |
| `WK_NOINTERNET` | `true/false` | (Optional) Whether to disable internet access, default is false. If public internet is not allowed, set it to true. Then mount risectl files to <risectl dir>/<version>/risectl. |
| `WK_RISECTLDIR` | `string` | (Optional) The path of the directory to store the risectl files, default is "$HOME/.risectl" |
| `WK_PROMETHEUS_DEFAULT_LABELS` | `string` | (Optional) The labels to be added to all queries. The format is "key1=value1,key2=value2". |


# Automated Initialization

WaveKit supports automated initialization through a configuration file, eliminating the need for manual setup through the web UI. This feature enables you to programmatically configure your WaveKit environment by pre-defining cluster and database information.

Here's an example initialization file structure:

```yaml
clusters:
  - name: Default Local Cluster
    version: v2.2.1
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
