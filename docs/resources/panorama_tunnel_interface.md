---
page_title: "panos: panos_panorama_tunnel_interface"
subcategory: "Panorama Networking"
---

# panos_panorama_tunnel_interface

This resource allows you to add/update/delete Panorama tunnel interfaces
for templates.


## Import Name

```
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_tunnel_interface" "example1" {
    template = panos_panorama_template.t.name
    name = "tunnel.5"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}

resource "panos_panorama_template" "t" {
    name = "foo"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  This must start with `tunnel.`.
* `template` - (Required) The template name.
* `vsys` - (Optional) The vsys that will use this interface (default: `vsys1`).
* `comment` - (Optional) The interface comment.
* `netflow_profile` - (Optional) The netflow profile.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
