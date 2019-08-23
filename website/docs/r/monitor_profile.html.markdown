---
layout: "panos"
page_title: "panos: panos_monitor_profile"
description: |-
  Manages monitor profiles.
---

# panos_monitor_profile

This resource allows you to add/update/delete monitor profiles.

**Minimum PAN-OS version**: 7.1

## Import Name

```
<name>
```

## Example Usage

```hcl
resource "panos_monitor_profile" "example" {
    name = "myProfile"
    interval = 5
    threshold = 3
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The monitor profile name.
* `interval` - (Optional, int) The probing interval in seconds.
* `threshold` - (Optional, int) The number of failed probes to determine that
  the tunnel is down.
* `action` - (Optional) Action triggered when tunnel's status changes.  Valid values
  are `wait-recover` (default) or `fail-over`.
