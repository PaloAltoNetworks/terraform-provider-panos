---
page_title: "panos: panos_ospf_area_interface"
subcategory: "Network"
---

# panos_ospf_area_interface

Gets info on an OSPF area interface.


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_area_interface" "example" {
    template = "my template"
    virtual_router = "my virtual router"
    ospf_area = "10.2.3.1"
    name = "ethernet1/1"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.
* `ospf_area` - (Required) OSPF area name.
* `name` - (Required) Interface name.


## Attribute Reference

The following attributes are supported:

* `enable` - (bool) Enable.
* `passive` - (bool) Passive.
* `link_type` - Link type.
* `metric` - (int) Metric.
* `priority` - (int) Priority.
* `hello_interval` - (int) Hello interval in seconds.
* `dead_counts` - (int) Dead counts.
* `retransmit_interval` - (int) Retransmit interval in seconds.
* `transit_delay` - (int) Transit delay in seconds.
* `grace_restart_delay` - (int) Graceful restart hello delay in seconds.
* `auth_profile` - Auth profile.
* `neighbors` - (list) List of neighbors.
* `bfd_profile` - BFD profile.
