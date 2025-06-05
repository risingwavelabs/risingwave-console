# Configuring RisingWave Console

## Configuration File

RisingWave Console is configured using a YAML file. This file includes several sections that define the settings for your RisingWave Console environment:

```yaml
init: string
host: string
port: integer
metricsport: integer
pg:
  dsn: string
root:
  password: string
nointernet: true/false
risectldir: string
worker:
  disable: true/false
ee:
  code: string
debug:
  enable: true/false
  port: integer

```

## Overriding Configuration with Environment Variables

You can override the YAML configuration settings by using environment variables. The following environment variables are supported:

| Environment Variable | Expected Value | Description |
|---------------------|----------------|-------------|
| `RCONSOLE_INIT` | `string` | (Optional) The path of file to store the initialization data, if not set, skip the initialization |
| `RCONSOLE_HOST` | `string` | (Optional) The host of the RisingWave Console server, it is used in the API endpoint of the web UI. If not set, the host will be localhost. |
| `RCONSOLE_PORT` | `integer` | (Optional) The port of the RisingWave Console server, default is 8020 |
| `RCONSOLE_METRICSPORT` | `integer` | (Optional) The port of the metrics server, default is 9020 |
| `RCONSOLE_PG_DSN` | `string` | (Required) The DSN (Data Source Name) for postgres database connection. If specified, Host, Port, User, Password, and Db settings will be ignored. |
| `RCONSOLE_ROOT_PASSWORD` | `string` | (Optional) The password of the root user, if not set, the default password is "123456" |
| `RCONSOLE_NOINTERNET` | `true/false` | (Optional) Whether to disable internet access, default is false. If public internet is not allowed, set it to true. Then mount risectl files to <risectl dir>/<version>/risectl. |
| `RCONSOLE_RISECTLDIR` | `string` | (Optional) The path of the directory to store the risectl files, default is "$HOME/.risectl" |
| `RCONSOLE_WORKER_DISABLE` | `true/false` | (Optional) Whether to disable the worker, default is false. |
| `RCONSOLE_EE_CODE` | `string` | (Optional) The activation code of the enterprise edition, if not set, the enterprise edition will be disabled. |
| `RCONSOLE_DEBUG_ENABLE` | `true/false` | (Optional) Whether to enable the debug server, default is false. |
| `RCONSOLE_DEBUG_PORT` | `integer` | (Optional) The port of the debug server, default is 8777 |


# Automated Initialization

RisingWave Console supports automated initialization through a configuration file, eliminating the need for manual setup through the web UI. This feature enables you to programmatically configure your RisingWave Console environment by pre-defining cluster and database information.

Here's an example initialization file structure:

```yaml
metricsStores:
  - name: Default Prometheus
    spec:
      prometheus:
        endpoint: http://prometheus:9500
clusters:
  - name: Default Local Cluster
    version: v2.2.1
    connections:
      host: rw
      sqlPort: 4566
      metaPort: 5690
      httpPort: 5691
    metricsStore: Default Prometheus
databases:
  - name: rw
    cluster: Default Local Cluster
    username: root
    database: dev

```

To use the initialization file, start RisingWave Console with the `RCONSOLE_INIT` environment variable pointing to your file:

```shell
RCONSOLE_INIT=/path/to/init.yaml risingwave-console
```



