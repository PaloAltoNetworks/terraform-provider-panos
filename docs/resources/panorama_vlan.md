---
page_title: "panos: panos_panorama_vlan"
subcategory: "Network"
---

# panos_panorama_vlan

This resource allows you to add/update/delete Panorama VLANs.


## PAN-OS

Panorama


## Import Name

```shell
<template>::<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_vlan" "example" {
    template = panos_panorama_template.t.name
    name = "myVlan"
    vlan_interface = panos_panorama_vlan_interface.vli.name

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_vlan_interface" "vli" {
    template = panos_panorama_template.t.name
    name = "vlan.6"
    vsys = "vsys1"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name.
* `vsys` - (Optional) The vsys to put the object into (default: `vsys1`).
* `template` - (Required) The template name.
* `vlan_interface` - (Optional) The VLAN interface.
* `interfaces` - (Optional, computed) List of layer2 interfaces.  You can also leave
  this blank and also use [panos_panorama_vlan_entry](./vlan_entry.html) for more control.
