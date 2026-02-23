# DNS Setup

Local DNS resolution using Pi-hole on `pihole` (192.168.1.2), with upstream DoH via Cloudflare.

## Pi-hole Installation

```bash
curl -sSL https://install.pi-hole.net | bash
```

During setup, select:
- Interface: `eth0`
- Upstream DNS: Cloudflare (1.1.1.1, 1.0.0.1)
- Install the web admin interface: yes

## Local DNS Records

Custom DNS entries configured in Pi-hole for accessing lab services by name:

| Hostname              | IP Address     | Service              |
|-----------------------|----------------|----------------------|
| pve.lab.local         | 192.168.1.10   | Proxmox web UI       |
| docker.lab.local      | 192.168.1.20   | Docker host           |
| nas.lab.local         | 192.168.1.30   | TrueNAS web UI       |
| grafana.lab.local     | 192.168.1.20   | Grafana dashboard    |
| pihole.lab.local      | 192.168.1.2    | Pi-hole admin        |
| wiki.lab.local        | 192.168.1.20   | This wiki            |

## DHCP DNS Configuration

OPNsense DHCP is configured to hand out `192.168.1.2` as the primary DNS server for all VLANs.

```
# OPNsense DHCP settings per VLAN
DNS Server 1: 192.168.1.2
DNS Server 2: 1.1.1.1  (fallback)
```

## Conditional Forwarding

To resolve local hostnames from Pi-hole, enable conditional forwarding:

```
# /etc/dnsmasq.d/02-local.conf
server=/lab.local/192.168.1.1
rev-server=192.168.1.0/24,192.168.1.1
```

See also: [Networking](networking.md) | [Docker Compose](docker_compose.md)
