# Old Server Rack Notes

These notes are from the original single-server setup before migrating to Proxmox.

## Hardware

- Dell PowerEdge T320
- 32 GB ECC RAM
- 4x 2TB SATA drives in RAID 10

## Services

Everything ran directly on Ubuntu Server 18.04:
- Nginx reverse proxy
- Nextcloud
- Plex Media Server
- Home Assistant

This setup was replaced in 2023 when we moved to a proper hypervisor-based architecture. See the current [Proxmox Setup](proxmox_setup.md) page.
