# Backup Strategy

Following the **3-2-1 rule**: 3 copies of data, on 2 different media, with 1 offsite.

## What Gets Backed Up

- Proxmox VM configs and disk images
- Docker volumes and compose files
- NAS datasets (documents, photos, media)
- Wiki `.storage/` directory
- OPNsense configuration XML

## Backup Layers

- **Layer 1 — Local snapshots:** ZFS snapshots on `pve-01` every 15 minutes, retained for 24 hours
- **Layer 2 — NAS replication:** Nightly ZFS send/receive from `pve-01` to `nas-01`
- **Layer 3 — Offsite:** Weekly encrypted backup to Hetzner StorageBox using restic

## Restic Commands

```bash
# Initialize a new restic repository
restic -r sftp:u123456@u123456.your-storagebox.de:/backups init

# Run a backup
restic -r sftp:u123456@u123456.your-storagebox.de:/backups \
  backup /mnt/data \
  --exclude="*.tmp" \
  --exclude=".cache"

# List snapshots
restic -r sftp:u123456@u123456.your-storagebox.de:/backups snapshots

# Prune old snapshots (keep 4 weekly, 6 monthly)
restic -r sftp:u123456@u123456.your-storagebox.de:/backups \
  forget --keep-weekly 4 --keep-monthly 6 --prune
```

## Automated Schedule

Backups are triggered by systemd timers on `docker-01`:

```bash
# /etc/systemd/system/restic-backup.timer
[Unit]
Description=Weekly restic backup

[Timer]
OnCalendar=Sun *-*-* 02:00:00
Persistent=true

[Install]
WantedBy=timers.target
```

See also: [Proxmox Setup](proxmox_setup.md)
