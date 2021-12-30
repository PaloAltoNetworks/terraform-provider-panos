---
page_title: "panos: panos_static_route_ipv4"
subcategory: "Network"
---

# panos_static_route_ipv4

This resource allows you to add/update/delete IPv4 static routes on a
virtual router.


## PAN-OS

NGFW


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_static_route_ipv4" "example" {
    name = "localnet"
    virtual_router = panos_virtual_router.vr1.name
    destination = "10.1.7.0/32"
    next_hop = "10.1.7.4"
}

resource "panos_virtual_router" "vr1" {
    name = "my virtual router"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The static route's name.
* `virtual_router` - (Required) The virtual router to add the static
  route to.
* `destination` - (Required) Destination IP address / prefix.
* `interface` - (Optional) Interface to use.
* `type` - (Optional) The next hop type.  Valid values are `ip-address` (the
  default), `discard`, `next-vr`, or an empty string for `None`.
* `next_hop` - (Optional) The value for the `type` setting.
* `admin_distance` - (Optional) The admin distance.
* `metric` - (Optional, int) Metric value / path cost (default: `10`).
* `route_table` - (Optional) Target routing table to install the route.  Valid
  values are `unicast` (the default), `no install`, `multicast`, or `both`.
* `bfd_profile` - (Optional, PAN-OS 7.1+) BFD configuration.
