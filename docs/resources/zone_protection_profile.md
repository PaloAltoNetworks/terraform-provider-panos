# panos_zone_protection_profile

Manages a Zone Protection Profile on a PAN-OS NGFW or Panorama.

Zone protection profiles define flood protection, packet-based attack protection, and reconnaissance protection settings that can be applied to security zones.

## Example Usage

### NGFW

```hcl
resource "panos_zone_protection_profile" "example" {
  location = {
    ngfw = {}
  }

  name        = "my-zone-protection"
  description = "Managed by Terraform"

  flood = {
    syn = {
      enable = true
      red = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
      syn_cookies = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
    }
    icmp = {
      enable = true
      red = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
    }
    icmpv6 = {
      enable = true
      red = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
    }
    udp = {
      enable = true
      red = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
    }
    other = {
      enable = true
      red = {
        alarm_rate    = 10000
        activate_rate = 20000
        maximal_rate  = 40000
      }
    }
  }

  discard_ip_spoof             = true
  discard_strict_source_routing = true
  discard_loose_source_routing  = true
  discard_malformed_option      = true
  remove_tcp_timestamp          = true
  discard_ip_frag               = false
}
```

### Panorama Template

```hcl
resource "panos_zone_protection_profile" "example" {
  location = {
    template = {
      name        = "my-template"
      ngfw_device = "localhost.localdomain"
    }
  }

  name        = "my-zone-protection"
  description = "Managed by Terraform"

  flood = {
    syn = {
      enable = true
      red = {
        alarm_rate    = 20000
        activate_rate = 25000
        maximal_rate  = 1000000
      }
    }
  }

  discard_ip_spoof             = true
  discard_strict_source_routing = true
  discard_loose_source_routing  = true
}
```

### Panorama Template Stack

```hcl
resource "panos_zone_protection_profile" "example" {
  location = {
    template_stack = {
      name        = "my-stack"
      ngfw_device = "localhost.localdomain"
    }
  }

  name = "my-zone-protection"
}
```

## Argument Reference

### location (Required)

Exactly one of the following location blocks must be specified:

- `ngfw` - Located directly on an NGFW device.
  - `ngfw_device` (Optional) - The NGFW device name. Defaults to `localhost.localdomain`.

- `template` - Located in a Panorama template.
  - `name` (Required) - The template name.
  - `panorama_device` (Optional) - The Panorama device. Defaults to `localhost.localdomain`.
  - `ngfw_device` (Optional) - The NGFW device. Defaults to `localhost.localdomain`.

- `template_stack` - Located in a Panorama template stack.
  - `name` (Required) - The template stack name.
  - `panorama_device` (Optional) - The Panorama device. Defaults to `localhost.localdomain`.
  - `ngfw_device` (Optional) - The NGFW device. Defaults to `localhost.localdomain`.

### Top-level Arguments

- `name` (Required) - The name of the zone protection profile.
- `description` (Optional) - A description for the profile.
- `discard_ip_spoof` (Optional) - Discard IP spoofed packets.
- `discard_strict_source_routing` (Optional) - Discard packets with the strict source routing IP option set.
- `discard_loose_source_routing` (Optional) - Discard packets with the loose source routing IP option set.
- `discard_malformed_option` (Optional) - Discard packets with malformed IP options.
- `remove_tcp_timestamp` (Optional) - Remove TCP timestamp option from packets.
- `discard_ip_frag` (Optional) - Discard IP fragments.
- `tcp_syn_with_data` (Optional) - Discard TCP SYN packets that carry data. Only supported on newer PAN-OS versions.
- `strip_tcp_fast_open_and_data` (Optional) - Strip TCP Fast Open option and data. Only supported on newer PAN-OS versions.
- `strip_mptcp_option` (Optional) - Strip MPTCP option. Valid values: `global`, `never`. Only supported on newer PAN-OS versions.

> **Note:** `tcp_syn_with_data`, `strip_tcp_fast_open_and_data`, and `strip_mptcp_option` are only accepted by PAN-OS versions that support the `<packet-based>` schema element. Do not set these fields if your Panorama/NGFW version does not support them — Panorama will reject the request with a schema validation error.

### flood (Optional)

Flood protection settings. Contains the following protocol blocks:

- `syn` - TCP SYN flood protection.
  - `enable` (Optional) - Enable SYN flood protection.
  - `red` (Optional) - Random Early Drop settings.
    - `alarm_rate` (Optional) - Alarm rate in packets/sec.
    - `activate_rate` (Optional) - Activate rate in packets/sec.
    - `maximal_rate` (Optional) - Maximum rate in packets/sec.
  - `syn_cookies` (Optional) - SYN Cookies settings (same attributes as `red`).

- `icmp` - ICMP flood protection.
  - `enable` (Optional) - Enable ICMP flood protection.
  - `red` (Optional) - Random Early Drop settings (same attributes as `syn.red`).

- `icmpv6` - ICMPv6 flood protection. Same structure as `icmp`.

- `udp` - UDP flood protection. Same structure as `icmp`.

- `other` - Other IP flood protection (`other-ip` in PAN-OS schema). Same structure as `icmp`.

## Dependency Ordering

If you are also managing a `panos_zone` resource that references this profile via `zone_protection_profile`, use a resource reference rather than a hardcoded string to ensure Terraform creates the profile before the zone:

```hcl
resource "panos_zone" "untrust" {
  # ...
  zone_protection_profile = panos_zone_protection_profile.example.name  # implicit dependency
}
```

If a hardcoded string is required, use `depends_on`:

```hcl
resource "panos_zone" "untrust" {
  zone_protection_profile = "my-zone-protection"
  depends_on = [panos_zone_protection_profile.example]
}
```

## Import

Zone protection profiles can be imported using a base64-encoded JSON import ID.

### NGFW

```shell
terraform import panos_zone_protection_profile.example \
  "$(echo -n '{"location":{"ngfw":{}},"name":"my-zone-protection"}' | base64)"
```

### Panorama Template

```shell
terraform import panos_zone_protection_profile.example \
  "$(echo -n '{"location":{"template":{"name":"my-template"}},"name":"my-zone-protection"}' | base64)"
```

### Panorama Template Stack

```shell
terraform import panos_zone_protection_profile.example \
  "$(echo -n '{"location":{"template_stack":{"name":"my-stack"}},"name":"my-zone-protection"}' | base64)"
```
