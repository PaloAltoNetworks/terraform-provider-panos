---
page_title: "panos: panos_panorama_vlan_entry"
subcategory: "Panorama Networking"
---

# panos_panorama_vlan_entry

This resource allows you to add/update/delete an interface in a VLAN on Panorama.


## Import Name

```
<template>::<vlan>:<interface>
```


## Example Usage

```hcl
resource "panos_panorama_vlan_entry" "example" {
    template = panos_panorama_template.t.name
    vlan = panos_vlan.vlan1.name
    interface = panos_ethernet_interface.e1.name
    mac_addresses = [
        "00:30:48:55:66:77",
        "00:30:48:55:66:88",
    ]
}

resource "panos_panorama_template" "t" {
    name = "my template"
}

resource "panos_panorama_vlan" "vlan1" {
    template = panos_panorama_template.t.name
    name = "myVlan"
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = panos_panorama_template.t.name
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"
}
```

## Argument Reference

The following arguments are supported:

* `vlan` - (Required) The VLAN's name.
* `template` - (Required) The template name.
* `interface` - (Required) The interface's name.
* `mac_addresses` - (Optional) List of MAC addresses that should go with this entry.
