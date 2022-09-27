---
page_title: "panos: panos_ospf_export"
subcategory: "Network"
---

# panos_ospf_export

Gets info on an OSPF export config attached to a virtual router.


## Import Name

NGFW:

```
<virtual_router>:<name>
```

Panorama:

```
<template>::<virtual_router>:<name>
```


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_export" "example" {
    template = panos_ospf_export.x.template
    virtual_router = panos_ospf_export.x.virtual_router
    name = panos_ospf_export.x.name
}

resource "panos_ospf_export" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = "10.2.3.0/24"
    tag = "10.5.15.151"
    metric = 42

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ospf" "x" {
    template = panos_panorama_template.x.name
    virtual_router = panos_panorama_virtual_router.x.name

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
* `name` - (Required) The export rule name.


## Attribute Reference

The following attributes are supported:

* `path_type` - Path type.
* `tag` - Tag.
* `metric` - (int) Metric.
