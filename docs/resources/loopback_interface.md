---
page_title: "panos: panos_loopback_interface"
subcategory: "Network"
---

# panos_loopback_interface

This resource allows you to add/update/delete loopback interfaces.


## PAN-OS

NGFW


## Import Name

```shell
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_loopback_interface" "example1" {
    name = "loopback.2"
    comment = "my loopback interface"
    static_ips = ["10.1.1.1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  This must start with `loopback.`.
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
