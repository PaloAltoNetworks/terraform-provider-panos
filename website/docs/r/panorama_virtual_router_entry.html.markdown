---
layout: "panos"
page_title: "panos: panos_panorama_virtual_router_entry"
sidebar_current: "docs-panos-panorama-resource-virtual-router-entry"
description: |-
  Manages an interface in a Panorama virtual router.
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
resource "panos_panorama_virtual_router" "vr" {
    template = "my template"
    name = "my vr"
}

resource "panos_panorama_virtual_router_entry" "example" {
    template = "my template"
    virtual_router = "${panos_panorama_virtual_router.vr.name}"
    interface = "ethernet1/5"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
* `virtual_router` - (Required) The virtual router's name.
* `interface` - (Required) The interface to import into the virtual router.
