# Proxmox VE Setup

Step-by-step notes for installing and configuring Proxmox VE 8.x on bare metal.

## Installation Steps

1. Download the latest Proxmox VE ISO from the official site
2. Flash to USB with `dd` or Balena Etcher
3. Boot from USB, select "Install Proxmox VE (Graphical)"
4. Choose the target disk (prefer ZFS mirror for redundancy)
5. Set hostname to `pve-01.lab.local` and configure a static IP
6. Complete installation and reboot

## Post-Install Configuration

Remove the enterprise repository and add the no-subscription repo:

```bash
# Disable enterprise repo
sed -i 's/^deb/#deb/' /etc/apt/sources.list.d/pve-enterprise.list

# Add no-subscription repo
echo "deb http://download.proxmox.com/debian/pve bookworm pve-no-subscription" \
  > /etc/apt/sources.list.d/pve-no-subscription.list

apt update && apt full-upgrade -y
```

## Enable IOMMU for PCI Passthrough

Edit the GRUB config:

```bash
# /etc/default/grub
GRUB_CMDLINE_LINUX_DEFAULT="quiet intel_iommu=on iommu=pt"
```

Then update GRUB and reboot:

```bash
update-grub
reboot
```

## Useful Commands

```bash
# List all VMs
qm list

# Start VM 100
qm start 100

# Check ZFS pool status
zpool status

# View cluster status
pvecm status
```

## Network Bridge Config

```
# /etc/network/interfaces
auto vmbr0
iface vmbr0 inet static
    address 192.168.1.10/24
    gateway 192.168.1.1
    bridge-ports enp3s0
    bridge-stp off
    bridge-fd 0
```

See also: [Docker in Proxmox](docker_compose.md) | [Backup Strategy](backup_strategy.md)
