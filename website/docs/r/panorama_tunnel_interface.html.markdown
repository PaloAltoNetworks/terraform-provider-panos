---
layout: "panos"
page_title: "panos: panos_panorama_tunnel_interface"
sidebar_current: "docs-panos-panorama-resource-tunnel-interface"
description: |-
  Manages Panorama tunnel interfaces.
---

# panos_panorama_tunnel_interface

This resource allows you to add/update/delete Panorama tunnel interfaces
for templates and template stacks.

## Example Usage

```hcl
resource "panos_panorama_tunnel_interface" "example1" {
    name = "tunnel.5"
    template_stack = "foo"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The interface's name.  This must start with `tunnel.`.
* `vsys` - (Optional) The vsys that will use this interface (default: `vsys1`).
* `comment` - (Optional) The interface comment.
* `netflow_profile` - (Optional) The netflow profile.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
