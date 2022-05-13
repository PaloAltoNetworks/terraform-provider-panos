---
page_title: "panos: panos_vlan"
subcategory: "Network"
---

# panos_vlan

This resource allows you to add/update/delete VLANs.


## PAN-OS

NGFW


## Import Name

```shell
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_vlan" "example" {
    name = "myVlan"
    vlan_interface = panos_vlan_interface.vli.name
}

resource "panos_vlan_interface" "vli" {
    name = "vlan.6"
    vsys = "vsys1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name.
* `vsys` - (Optional) The vsys to put the object into (default: `vsys1`).
* `vlan_interface` - (Optional) The VLAN interface.
* `interfaces` - (Optional, computed) List of layer2 interfaces.  You can also leave
  this blank and also use [panos_vlan_entry](vlan_entry.html) for more control.
