# Home Network Overview

Network topology and VLAN configuration for the lab, managed by OPNsense on `gw-01`.

## VLAN Layout

| VLAN ID | Name       | Subnet           | Purpose                  |
|---------|------------|------------------|--------------------------|
| 1       | Default    | 192.168.1.0/24   | Management & servers     |
| 10      | Trusted    | 192.168.10.0/24  | Personal devices         |
| 20      | IoT        | 192.168.20.0/24  | Smart home devices       |
| 30      | Guest      | 192.168.30.0/24  | Guest Wi-Fi              |
| 99      | DMZ        | 192.168.99.0/24  | Publicly exposed services|

## Firewall Rules Summary

Key inter-VLAN rules on OPNsense:

```
# Allow Trusted -> Management (for admin access)
pass  Trusted_net  ->  Management_net  any

# Allow Trusted -> IoT (for controlling devices)
pass  Trusted_net  ->  IoT_net         any

# Block IoT -> all other VLANs
block IoT_net      ->  Trusted_net     any
block IoT_net      ->  Management_net  any

# Allow Guest -> WAN only (internet access)
pass  Guest_net    ->  WAN             any
block Guest_net    ->  RFC1918         any
```

> **Tip:** Always place more specific rules above general block rules in OPNsense. The firewall evaluates rules top-down and stops at the first match.

## Physical Topology

```
ISP Modem
  └── OPNsense (gw-01)
        ├── UniFi Switch 24 (trunk)
        │     ├── pve-01 (VLAN 1)
        │     ├── docker-01 (VLAN 1)
        │     ├── nas-01 (VLAN 1)
        │     ├── UniFi AP (VLANs 10, 20, 30)
        │     └── IoT devices (VLAN 20)
        └── pihole (VLAN 1)
```

See also: [DNS Setup](dns_setup.md) | [Docker Compose](docker_compose.md)
