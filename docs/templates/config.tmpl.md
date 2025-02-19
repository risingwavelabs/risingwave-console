# Configuring WaveKit

## Configuration File

WaveKit is configured using a YAML file. This file includes several sections that define the settings for your WaveKit environment:

```yaml
{{CONFIG_SAMPLE_YAML}}
```

## Overriding Configuration with Environment Variables

You can override the YAML configuration settings by using environment variables. The following environment variables are supported:

{{CONFIG_ENV}}

# Automated Initialization

WaveKit supports automated initialization through a configuration file, eliminating the need for manual setup through the web UI. This feature enables you to programmatically configure your WaveKit environment by pre-defining cluster and database information.

Here's an example initialization file structure:

```yaml
{{CONFIG_SAMPLE_INIT}}
```

To use the initialization file, start WaveKit with the `WK_INIT` environment variable pointing to your file:

```shell
WK_INIT=/path/to/init.yaml wavekit
```
