# RisingWave Console

**An On-Premise Observability and Management UI for RisingWave**

![RisingWave Console Cover](docs/images/cover.png)

RisingWave Console is a user-friendly, on-premise tool designed to enhance observability and simplify management for your existing RisingWave clusters. It provides a web-based interface to connect to, monitor, and interact with RisingWave, whether it's deployed via Docker, Kubernetes, or on RisingWave Cloud.

## ✨ Key Features

*   **Centralized Cluster View:** Connect to and manage all your RisingWave instances from a single dashboard.
*   **Status Monitoring:** Get at-a-glance insights into the status and basic information of your connected clusters.
*   **Direct `risectl` Execution:** Run `risectl` commands on your clusters directly from the RisingWave Console UI.
*   **Metadata Snapshot Management:** Easily create manual snapshots and configure automated backups for your RisingWave cluster metadata.
*   **Diagnostic Collection:** Trigger and automate the collection of diagnostic information for advanced troubleshooting.
*   **Interactive SQL Console:**
    *   Securely connect to specific databases within your RisingWave clusters.
    *   Explore schemas, tables, views, materialized views, sources, and sinks.
    *   Execute SQL queries with multi-tab support and query history.
    *   Visualize data flows with an integrated **Streaming Graph** for streaming queries.
*   **On-Premise Control:** Host RisingWave Console within your own environment, ensuring your connection details and metadata remain under your control.

## 📋 Prerequisites

*   **A running RisingWave cluster:** RisingWave Console connects to existing, operational RisingWave instances.
*   **PostgreSQL Database for RisingWave Console Metadata:** RisingWave Console uses a PostgreSQL database to store its own configuration (e.g., connection details for your RisingWave clusters). This can be self-hosted or the one bundled with specific RisingWave Console Docker images.
*   **Docker (Recommended):** If you plan to use the Docker-based installation methods for RisingWave Console.

## 🚀 Quick Start (Docker with Persistent Storage)

This is the recommended way to get started quickly with RisingWave Console, ensuring your RisingWave Console metadata persists.

1.  Ensure Docker is installed and running.
2.  Run the following command (replace `vX.Y.Z` with the latest RisingWave Console version, e.g., `v0.4.0`):
    ```shell
    docker run -d -p 8020:8020 --name risingwave-console \
      -e RCONSOLE_ROOT_PASSWORD=your_secure_password \
      -v risingwave-console-data:/var/lib/postgresql \
      risingwavelabs/risingwave-console:vX.Y.Z-pgbundle
    ```
    *   This uses the `-pgbundle` image which includes PostgreSQL.
    *   `-d` runs the container in detached mode.
    *   `-p 8020:8020` maps the port.
    *   `-e RCONSOLE_ROOT_PASSWORD=your_secure_password` sets a custom initial password for the `root` user. **Recommended!**
    *   `-v risingwave-console-data:/var/lib/postgresql` creates a Docker volume named `risingwave-console-data` to persist PostgreSQL data.

3.  Access RisingWave Console: Open your browser and go to `http://localhost:8020`.
4.  Login with:
    *   **Username:** `root`
    *   **Password:** `your_secure_password` (or `root` if `RCONSOLE_ROOT_PASSWORD` was not set).

## 🛠️ Other Installation Methods

RisingWave Console offers flexibility in how it can be deployed:

*   **Docker (Ephemeral Storage):** For quick testing without persistence. Metadata is lost if the container is removed.
    ```shell
    docker run --rm -p 8020:8020 --name risingwave-console risingwavelabs/risingwave-console:vX.Y.Z-pgbundle
    ```
*   **Standalone Binary:** Download the binary and run it directly, connecting to your own existing PostgreSQL instance.
    ```shell
    # 1. Download
    curl https://risingwave-console.s3.ap-southeast-1.amazonaws.com/download.sh | sh
    # 2. Run (ensure RCONSOLE_PG_DSN is set)
    RCONSOLE_PG_DSN="postgres://user:pass@host:port/dbname" RCONSOLE_ROOT_PASSWORD=your_secure_password ./risingwave-console
    ```
*   **Docker Compose (Recommended for Production with Self-Managed PG):** Use Docker Compose to manage RisingWave Console and (optionally) a dedicated PostgreSQL container. See example `docker-compose.yaml` in our [Installation Guide](risingwave.com/wavekit/installation-setup).
    ```yaml
    # Example snippet for docker-compose.yaml using non-pgbundle image
    # services:
    #   risingwave-console:
    #     image: risingwavelabs/risingwave-console:vX.Y.Z
    #     ports: ["8020:8020"]
    #     environment:
    #       RCONSOLE_PG_DSN: "your_postgres_dsn"
    #       RCONSOLE_ROOT_PASSWORD: "your_secure_password"
    #   # ... your PostgreSQL service definition ...
    ```

## 💻 Using RisingWave Console

Once RisingWave Console is running and you've logged in:

1.  **Connect Your RisingWave Cluster:**
    *   Navigate to the "Clusters" section.
    *   Click "Add Cluster" and provide the connection details (Host, SQL Port, Meta Node Port, HTTP Port, Version) for your existing RisingWave cluster.
    *   Test and save the connection.

2.  **Explore Cluster Details:**
    *   From the "Clusters" list, click on your connected cluster to view its details page.
    *   Here you can:
        *   View cluster information and status.
        *   Execute `risectl` commands.
        *   Manage metadata snapshots.
        *   Collect diagnostic information.

3.  **Use the SQL Console:**
    *   Navigate to the "SQL Console" section.
    *   Click "Manage Databases" to add a new database connection, linking it to one of your configured clusters and providing RisingWave database credentials.
    *   Write and execute SQL queries, explore schemas, view query history, and visualize streaming graphs.

## ⚙️ Configuration

RisingWave Console can be configured using environment variables. For a detailed list of available settings, please refer to our [Configuration Guide](docs/config.md).

Key variables include:
*   `RCONSOLE_PORT`: Port for the RisingWave Console UI (default: `8020`).
*   `RCONSOLE_PG_DSN`: PostgreSQL connection string for RisingWave Console's metadata.
*   `RCONSOLE_ROOT_PASSWORD`: Initial password for the `root` UI user (default: `root`).
*   `RCONSOLE_RISECTLDIR`: Path related to `risectl` resources if needed.

## 📚 Documentation

For comprehensive information, guides, and usage details, please visit our **[Official Documentation](risingwave.com/wavekit/introduction)**

##  editions

RisingWave Console is available in two editions:

*   **RisingWave Console-Lite:** The open-source edition (Apache 2.0) with core functionalities. This is what you get by default.
*   **RisingWave Console-Pro:** A future edition with advanced capabilities for enterprise users. A license key will be required. To apply or express interest, contact us at [sales@risingwave-labs.com](mailto:sales@risingwave-labs.com) or [fill out this form](https://cloud.risingwave.com/auth/license_key/).

## 🤝 Contributing to RisingWave Console

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to the project, report bugs, or request features.

## 📄 License

RisingWave Console-Lite is licensed under the [Apache License 2.0](LICENSE).
