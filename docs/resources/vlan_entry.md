---
page_title: "panos: panos_vlan_entry"
subcategory: "Network"
---

# panos_vlan_entry

This resource allows you to add/update/delete an interface in a VLAN.


## PAN-OS

NGFW


## Import Name

```shell
<vlan>:<interface>
```


## Example Usage

```hcl
resource "panos_vlan_entry" "example" {
    vlan = panos_vlan.vlan1.name
    interface = panos_ethernet_interface.e1.name
    mac_addresses = [
        "00:30:48:55:66:77",
        "00:30:48:55:66:88",
    ]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_vlan" "vlan1" {
    name = "myVlan"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ethernet_interface" "e1" {
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `vlan` - (Required) The VLAN's name.
* `interface` - (Required) The interface's name.
* `mac_addresses` - (Optional) List of MAC addresses that should go with this entry.
