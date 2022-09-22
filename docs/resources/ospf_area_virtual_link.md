---
page_title: "panos: panos_ospf_area_virtual_link"
subcategory: "Network"
---

# panos_ospf_area_virtual_link

Manages an OSPF area virtual link.


## Import Name

NGFW:

```shell
<virtual_router>:<ospf_area>:<name>
```

Panorama:

```shell
<template>::<virtual_router>:<ospf_area>:<name>
```


## Example Usage

```hcl
# Panorama example.
resource "panos_ospf_area_virtual_link" "example" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = "foo"
    neighbor_id = "10.20.40.80"
    transit_area_id = panos_ospf_area.x.name

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ospf_area" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = "10.30.40.50"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = "vr name here"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "x" {
    name = "my template"

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.
* `ospf_area` - (Required) OSPF area name.
* `name` - (Required) Interface name.
* `enable` - (bool) Enable (default: `true`).
* `neighbor_id` - (Required) Neighbor ID.
* `transit_area_id` - (Required) Transit area ID.
* `hello_interval` - (int) Hello interval in seconds (default: `10`).
* `dead_counts` - (int) Dead counts (default: `4`).
* `retransmit_interval` - (int) Retransmit interval in seconds (default: `5`).
* `transit_delay` - (int) Transit delay in seconds (default: `1`).
* `auth_profile` - Auth profile.
* `bfd_profile` - BFD profile.
