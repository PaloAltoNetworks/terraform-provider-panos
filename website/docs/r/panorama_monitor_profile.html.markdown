---
layout: "panos"
page_title: "panos: panos_panorama_monitor_profile"
description: |-
  Manages Panorama monitor profiles.
---

# panos_panorama_monitor_profile

This resource allows you to add/update/delete Panorama monitor profiles.

**Minimum PAN-OS version**: 7.1

## Import Name

```
<template>:<template_stack>:<name>
```

## Example Usage

```hcl
resource "panos_panorama_monitor_profile" "example" {
    name = "myProfile"
    interval = 5
    threshold = 3
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The monitor profile name.
* `interval` - (Optional, int) The probing interval in seconds.
* `threshold` - (Optional, int) The number of failed probes to determine that
  the tunnel is down.
* `action` - (Optional) Action triggered when tunnel's status changes.  Valid values
  are `wait-recover` (default) or `fail-over`.
