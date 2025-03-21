mkdir -p promduck/dashboards

URL_BASE="https://raw.githubusercontent.com/risingwavelabs/wavekit-lite/refs/heads/main/examples/promduck"

curl -L $URL_BASE/docker-compose.yaml -o promduck/docker-compose.yaml
curl -L $URL_BASE/grafana-risedev-dashboard.yml -o promduck/grafana-risedev-dashboard.yml
curl -L $URL_BASE/grafana-risedev-datasource.yml -o promduck/grafana-risedev-datasource.yml
curl -L $URL_BASE/grafana.ini -o promduck/grafana.ini
curl -L $URL_BASE/prometheus.yaml -o promduck/prometheus.yaml
curl -L https://raw.githubusercontent.com/risingwavelabs/risingwave/refs/heads/main/docker/dashboards/risingwave-dev-dashboard.json -o promduck/dashboards/risingwave-dev-dashboard.json
