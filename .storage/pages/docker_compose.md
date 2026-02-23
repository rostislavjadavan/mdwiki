# Docker Compose Recipes

Useful Docker Compose stacks running on `docker-01`.

## Traefik Reverse Proxy

Traefik handles TLS termination and routes traffic to containers automatically using labels.

```yaml
services:
  traefik:
    image: traefik:v3.1
    command:
      - "--api.dashboard=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.le.acme.httpchallenge.entrypoint=web"
      - "--certificatesresolvers.le.acme.email=admin@lab.local"
      - "--certificatesresolvers.le.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - traefik-certs:/letsencrypt

volumes:
  traefik-certs:
```

## Monitoring Stack (Grafana + Prometheus)

A basic monitoring setup with Prometheus scraping node metrics and Grafana for dashboards.

```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    ports:
      - "9090:9090"

  grafana:
    image: grafana/grafana:latest
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=changeme
    volumes:
      - grafana-data:/var/lib/grafana
    ports:
      - "3000:3000"

  node-exporter:
    image: prom/node-exporter:latest
    pid: host
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
    command:
      - "--path.procfs=/host/proc"
      - "--path.sysfs=/host/sys"

volumes:
  prometheus-data:
  grafana-data:
```

## Managing Stacks

```bash
# Start a stack in detached mode
docker compose up -d

# View logs for a specific service
docker compose logs -f grafana

# Restart a single service
docker compose restart traefik

# Pull updated images and recreate
docker compose pull && docker compose up -d
```

See also: [Proxmox Setup](proxmox_setup.md) | [Networking](networking.md)
