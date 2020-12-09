---
page_title: "panos: panos_panorama_vlan"
subcategory: "Panorama Networking"
---

# panos_panorama_vlan

This resource allows you to add/update/delete Panorama VLANs.


## Import Name

```
<template>::<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_vlan" "example" {
    template = panos_panorama_template.t.name
    name = "myVlan"
    vlan_interface = panos_panorama_vlan_interface.vli.name
}

resource "panos_panorama_vlan_interface" "vli" {
    template = panos_panorama_template.t.name
    name = "vlan.6"
    vsys = "vsys1"
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name.
* `vsys` - (Optional) The vsys to put the object into (default: `vsys1`).
* `template` - (Required) The template name.
* `vlan_interface` - (Optional) The VLAN interface.
* `interfaces` - (Optional, computed) List of layer2 interfaces.  You can also leave
  this blank and also use [panos_vlan_entry](./vlan_entry.html) for more control.
