# WaveKit

WaveKit is a simple on-prem tool designed to enhance observability for your RisingWave cluster, enabling faster issue detection, efficient troubleshooting, and improved performance.

![WaveKit Cover](docs/images/cover.png)

WaveKit supports all RisingWave deployment types, including Docker, Kubernetes, and RisingWave Cloud.

> [!WARNING]
> **WaveKit is currently in the public preview stage.**

> [!NOTE]
> _WaveKit uses a PostgreSQL database to store key cluster metadata, including connection details like hostnames and ports for RisingWave clusters. To ensure persistence, you’ll need to self-host a PostgreSQL database to prevent metadata loss._

> [!NOTE]
> _To use WaveKit, ensure your RisingWave cluster is already running and accessible._


## Installation (Quick setup with Docker)

This method installs WaveKit with a bundled PostgreSQL database for convenience. However, if you prefer to use your own self-hosted PostgreSQL database for data persistence, skip to the next section.  

### **Starting the WaveKit Server**  

You can start the WaveKit server in two ways:  

#### **Option 1: Ephemeral Storage (No Persistence)**  
Runs WaveKit with a bundled PostgreSQL database, but metadata is stored inside the container. If the container is removed, all metadata will be lost.  

```shell
docker run --rm -p 8020:8020 --name wavekit risingwavelabs/wavekit:v0.1.2-pgbundle
```

#### **Option 2: Persistent Storage (Recommended)**  
Runs WaveKit with a bundled PostgreSQL database and stores metadata in a persistent Docker volume (`wavekit-data`), ensuring data persists across restarts.  

```shell
docker run -p 8020:8020 --name wavekit -v wavekit-data:/var/lib/postgresql risingwavelabs/wavekit:v0.1.2-pgbundle
```

### **Accessing WaveKit**  

Once the server is running, open your browser and go to:  

- **[http://localhost:8020](http://localhost:8020)**  

Use the following default credentials to log in:  
- **Username:** `root`  
- **Password:** `root`  


## Installation (Quick setup with Binary)

To install WaveKit using a standalone binary, follow these steps:

1. Download and install the latest WaveKit binary:

  ```shell
  curl https://wavekit-release.s3.ap-southeast-1.amazonaws.com/download.sh | sh
  ```

2. Run the following command to start the WaveKit server:

  ```shell
  WK_PG_DSN=postgres://postgres:postgres@localhost:5432/postgres ./wavekit
  ```

  > [!NOTE]
  > Ensure you have a PostgreSQL database running on your machine and set the WK_PG_DSN environment variable to your database connection string.

## Installation (Recommended for production)

The following section provides a step-by-step guide to setting up WaveKit with your self-hosted PostgreSQL database. This approach is recommended if you need persistent metadata with high availability.

First, create `docker-compose.yaml` file with the following content:

```yaml
version: "3.9"
services:
  wavekit:
    image: cloudcarver/wavekit:v0.1.2
    ports:
      - "8020:8020"
    environment:
      WK_PORT: 8020
      WK_PG_DSN: postgres://postgres:postgres@localhost:5432/postgres
      WK_RISECTLDIR: /

  db: 
    image: "postgres:latest"
    ports:
      - "5432:5432"
    environment:
      POSTGRES_PASSWORD: postgres
    volumes:
      - db-data:/var/lib/postgresql/data

volumes:
  db-data:
  wavekit-data:
```

Start WaveKit by running the following command:

```shell
docker compose up
```

## Customizing WaveKit Settings

WaveKit offers flexible configuration options through either a configuration file or environment variables. For detailed information about available settings and configuration methods, please refer to our [configuration documentation](docs/config.md).


## Contributing to WaveKit

We welcome contributions to WaveKit! Please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) for more information on how to contribute to the project.


## WaveKit Editions

WaveKit is available in two editions:  

- **WaveKit-Lite** – A lightweight, open-source edition that includes core functionalities. Licensed under Apache 2.0.  
- **WaveKit-Pro** – A full-featured edition with advanced capabilities. A license key is required for access. To apply, contact us at [sales@risingwave-labs.com](mailto:sales@risingwave-labs.com) or [fill out this form](https://cloud.risingwave.com/auth/license_key/).
