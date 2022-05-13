---
page_title: "panos: panos_virtual_router_entry"
subcategory: "Network"
---

# panos_virtual_router_entry

This resource allows you to add/update/delete an interface in a
virtual router.

This resource has some overlap with the `panos_virtual_router`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_virtual_router` spec does not define the
`interfaces` field.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_virtual_router_entry`


## Import Name

```shell
<template>:<template_stack>:<virtual_router>:<interface>
```


## Example Usage

```hcl
resource "panos_virtual_router_entry" "example" {
    virtual_router = panos_virtual_router.vr.name
    interface = panos_ethernet_interface.e.name
}

resource "panos_virtual_router" "vr" {
    name = "my vr"
}

resource "panos_ethernet_interface" "e" {
    name = "ethernet1/1"
    mode = "layer3"
}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


The following arguments are supported:

* `virtual_router` - (Required) The virtual router's name.
* `interface` - (Required) The interface to import into the virtual router.
