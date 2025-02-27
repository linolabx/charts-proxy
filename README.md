# Charts Proxy

## Usage

```bash
mkdir charts
cat <<EOF > config.yaml
target_dir: /charts

default_proxy: http://your-proxy-server:port

repos:
  - name: keel
    kind: HELM
    url: https://charts.keel.sh
    cron: "0 * * * *"
    charts:
      - name: keel
        versions_from: 1.0.3
  - name: cert-manager
    kind: HELM
    url: https://charts.jetstack.io
    cron: "0 * * * *"
    charts:
      - name: cert-manager
        versions_from: 1.13.0
        stable_only: true
EOF

docker run -d --name charts-proxy \
 -e CONFIG_FILE=/config.yaml \
 -v ./config.yaml:/config.yaml \
 -v ./charts:/charts \
 -p 8080:8080 \
 ghcr.io/linolabx/charts-proxy:latest


# host /charts by your favorite web server
```

# TODO:

- [ ] internal web server
