---
page_title: "panos: panos_ospf_area_virtual_link"
subcategory: "Network"
---

# panos_ospf_area_virtual_link

Gets info on an OSPF area virtual link.


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_area_virtual_link" "example" {
    template = "my template"
    virtual_router = "the virtual router name"
    ospf_area = "10.10.10.10"
    name = "foo"
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
* `neighbor_id` - Neighbor ID.
* `transit_area_id` - Transit area ID.
* `hello_interval` - (int) Hello interval in seconds.
* `dead_counts` - (int) Dead counts.
* `retransmit_interval` - (int) Retransmit interval in seconds.
* `transit_delay` - (int) Transit delay in seconds.
* `auth_profile` - Auth profile.
* `bfd_profile` - BFD profile.
