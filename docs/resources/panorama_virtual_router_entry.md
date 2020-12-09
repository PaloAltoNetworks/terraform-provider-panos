---
page_title: "panos: panos_panorama_virtual_router_entry"
subcategory: "Panorama Networking"
---

# panos_panorama_virtual_router_entry

This resource allows you to add/update/delete an interface in a Panorama
virtual router template.

This resource has some overlap with the `panos_panorama_virtual_router`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_virtual_router` spec does not define the
`interfaces` field.


## Import Name

```
<template>:<template_stack>:<virtual_router>:<interface>
```


## Example Usage

```hcl
resource "panos_panorama_virtual_router_entry" "example" {
    template = panos_panorama_template.tmpl.name
    virtual_router = panos_panorama_virtual_router.vr.name
    interface = panos_panorama_ethernet_interface.e1.name
}

resource "panos_panorama_template" "tmpl" {
    name = "my template"
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.tmpl.name
    name = "my vr"
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = panos_panorama_template.tmpl.name
    name = "ethernet1/5"
    mode = "layer3"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
* `virtual_router` - (Required) The virtual router's name.
* `interface` - (Required) The interface to import into the virtual router.
