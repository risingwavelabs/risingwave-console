mkdir -p prompush/dashboards

URL_BASE="https://raw.githubusercontent.com/risingwavelabs/wavekit-lite/refs/heads/main/examples/prompush"

curl -L $URL_BASE/docker-compose.yaml -o prompush/docker-compose.yaml
curl -L $URL_BASE/grafana-risedev-dashboard.yml -o prompush/grafana-risedev-dashboard.yml
curl -L $URL_BASE/grafana-risedev-datasource.yml -o prompush/grafana-risedev-datasource.yml
curl -L $URL_BASE/grafana.ini -o prompush/grafana.ini
curl -L $URL_BASE/prometheus.yaml -o prompush/prometheus.yaml
curl -L https://raw.githubusercontent.com/risingwavelabs/risingwave/refs/heads/main/docker/dashboards/risingwave-dev-dashboard.json -o prompush/dashboards/risingwave-dev-dashboard.json
