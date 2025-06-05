# RisingWave Console

RisingWave Console (previously known as WaveKit) is a simple on-prem tool designed to enhance observability for your RisingWave cluster, enabling faster issue detection, efficient troubleshooting, and improved performance.

![RisingWave Console Cover](docs/images/cover.png)

RisingWave Console supports all RisingWave deployment types, including Docker, Kubernetes, and RisingWave Cloud.

> [!NOTE]
> _RisingWave Console uses a PostgreSQL database to store key cluster metadata, including connection details like hostnames and ports for RisingWave clusters. To ensure persistence, you’ll need to self-host a PostgreSQL database to prevent metadata loss._

> [!NOTE]
> _To use RisingWave Console, ensure your RisingWave cluster is already running and accessible._


## Installation (Quick setup with Docker)

This method installs RisingWave Console with a bundled PostgreSQL database for convenience. However, if you prefer to use your own self-hosted PostgreSQL database for data persistence, skip to the next section.  

### **Starting the RisingWave Console Server**  

You can start the RisingWave Console server in two ways:  

#### **Option 1: Ephemeral Storage (No Persistence)**  
Runs RisingWave Console with a bundled PostgreSQL database, but metadata is stored inside the container. If the container is removed, all metadata will be lost.  

```shell
docker run --rm -p 8020:8020 --name risingwave-console risingwavelabs/risingwave-console:v0.4.0-pgbundle
```

#### **Option 2: Persistent Storage (Recommended)**  
Runs RisingWave Console with a bundled PostgreSQL database and stores metadata in a persistent Docker volume (`console-data`), ensuring data persists across restarts.  

```shell
docker run -p 8020:8020 --name risingwave-console -v console-data:/var/lib/postgresql risingwavelabs/risingwave-console:v0.4.0-pgbundle
```

### **Accessing RisingWave Console**  

Once the server is running, open your browser and go to:  

- **[http://localhost:8020](http://localhost:8020)**  

Use the following default credentials to log in:  
- **Username:** `root`  
- **Password:** `root`  


## Installation (Quick setup with Binary)

To install RisingWave Console using a standalone binary, follow these steps:

1. Download and install the latest RisingWave Console binary:

  ```shell
  curl https://risingwave-console.s3.ap-southeast-1.amazonaws.com/download.sh | sh
  ```

2. Run the following command to start the RisingWave Console server:

  ```shell
  WK_PG_DSN=postgres://postgres:postgres@localhost:5432/postgres WK_ROOT_PASSWORD=root ./risingwave-console
  ```

  > [!NOTE]
  > Ensure you have a PostgreSQL database running on your machine and set the WK_PG_DSN environment variable to your database connection string.

## Installation (Recommended for production)

The following section provides a step-by-step guide to setting up RisingWave Console with your self-hosted PostgreSQL database. This approach is recommended if you need persistent metadata with high availability.

First, create `docker-compose.yaml` file with the following content:

```yaml
version: "3.9"
services:
  console:
    image: risingwavelabs/risingwave-console:v0.4.0
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
```

Start RisingWave Console by running the following command:

```shell
docker compose up
```

## Customizing RisingWave Console Settings

RisingWave Console offers flexible configuration options through either a configuration file or environment variables. For detailed information about available settings and configuration methods, please refer to our [configuration documentation](docs/config.md).

## RisingWave Console Editions

RisingWave Console is available in two editions:  

- **RisingWave Console Lite** – A lightweight, open-source edition that includes core functionalities. Licensed under Apache 2.0.  
- **RisingWave Console Pro** – A full-featured edition with advanced capabilities. A license key is required for access. To apply, contact us at [sales@risingwave-labs.com](mailto:sales@risingwave-labs.com) or [fill out this form](https://cloud.risingwave.com/auth/license_key/).


## Contributing to RisingWave Console

We welcome contributions to RisingWave Console! Please refer to our [CONTRIBUTING.md](CONTRIBUTING.md) for more information on how to contribute to the project.
