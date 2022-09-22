---
page_title: "panos: panos_panorama_loopback_interface"
subcategory: "Network"
---

# panos_panorama_loopback_interface

This resource allows you to add/update/delete Panorama loopback interfaces
for templates.


## PAN-OS

Panorama


## Import Name

```shell
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_loopback_interface" "example" {
    name = "loopback.2"
    template = panos_panorama_template.t.name
    comment = "my loopback interface"
    static_ips = ["10.1.1.1"]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  This must start with `loopback.`.
* `template` - (Required) The template name.
* `vsys` - (Optional) The vsys that will use this interface (default: `vsys1`).
* `comment` - (Optional) The interface comment.
* `netflow_profile` - (Optional) The netflow profile.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
* `adjust_tcp_mss` - (Optional, bool) Adjust TCP MSS (default: false).
* `ipv4_mss_adjust` - (Optional, PAN-OS 8.0+) The IPv4 MSS adjust value.
* `ipv6_mss_adjust` - (Optional, PAN-OS 8.0+) The IPv6 MSS adjust value.
