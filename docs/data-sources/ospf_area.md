---
page_title: "panos: panos_ospf_area"
subcategory: "Network"
---

# panos_ospf_area

Gets info on an OSPF area.


## Example Usage

```hcl
data "panos_ospf_area" "example" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    name = panos_ospf_area.x.name
}

resource "panos_ospf_area" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = "10.2.3.4"
    type = "nssa"
    accept_summary = true
    default_route_advertise = true
    advertise_metric = 50
    advertise_type = "ext-2"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ospf" "x" {
    template = panos_panorama_template.x.name
    virtual_router = panos_panorama_virtual_router.x.name
    enable = true

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = "my virtual router"

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
* `name` - (Required) Name.

## Attribute Reference

The following attributes are available:

* `type` - Area type.
* `accept_summary` - (bool) (stub/nssa) Accept summary.
* `default_route_advertise` - (bool) (stub/nssa) Default route advertise.
* `advertise_metric` - (int) (stub/nssa) Advertise metric.
* `advertise_type` - (nssa) Advertise type.
* `ext_range` - (repeatable) (nssa) EXT range spec, as defined below.
* `range` - (repeatable) Range spec, as defined below.

`ext_range` and `range` both support the following arguments:

* `network` - (Required) Network.
* `action` - Action.
