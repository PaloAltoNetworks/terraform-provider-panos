---
page_title: "panos: panos_panorama_layer2_subinterface"
subcategory: "Panorama Networking"
---

# panos_panorama_layer2_subinterface

This resource allows you to add/update/delete Panorama layer2 subinterfaces.


## Import Name

```
<template>::<interface_type>:<parent_interface>:<parent_mode>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_layer2_subinterface" "example" {
    template = panos_panorama_template.tmpl.name
    parent_interface = panos_panorama_ethernet_interface.e.name
    parent_mode = panos_panorama_ethernet_interface.e.mode
    vsys = "vsys1"
    name = "ethernet1/5.5"
    tag = 5
}

resource "panos_panorama_ethernet_interface" "e" {
    template = panos_panorama_template.tmpl.name
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer2"
}

resource "panos_panorama_template" "tmpl" {
    name = "myTemplate"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
* `interface_type` - (Optional) The interface type.  Valid values are `ethernet` (default)
  or `aggregate-ethernet`.
* `parent_interface` - (Required) The name of the parent interface.
* `parent_mode` - (Optional) The parent's mode.  Valid values are `layer2` (default)
  or `virtual-wire`.
* `vsys` - (Required) The vsys that will use this interface.  This should be
  something like `vsys1` or `vsys3`.
* `name` - (Required) The interface's name.
* `tag` - (Optional, int) The interface's tag.
* `netflow_profile` - (Optional) The netflow profile.
* `comment` - (Optional) The interface comment.
