---
layout: "panos"
page_title: "panos: panos_tunnel_interface"
sidebar_current: "docs-panos-resource-tunnel-interface"
description: |-
  Manages tunnel interfaces.
---

# panos_tunnel_interface

This resource allows you to add/update/delete tunnel interfaces.

## Example Usage

```hcl
resource "panos_tunnel_interface" "example1" {
    name = "tunnel.5"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  This must start with `tunnel.`.
* `vsys` - (Optional) The vsys that will use this interface (default: `vsys1`).
* `comment` - (Optional) The interface comment.
* `netflow_profile` - (Optional) The netflow profile.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
