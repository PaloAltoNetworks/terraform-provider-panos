---
page_title: "panos: panos_arp"
subcategory: "Network"
---

# panos_arp

Gets ARP config attached to an interface.


## Example Usage

```hcl
# Panorama ethernet interface example.
data "panos_arp" "example1" {
    template = panos_arp.x.template
    interface_type = panos_arp.x.interface_type
    interface_name = panos_arp.x.interface_name
    ip = panos_arp.x.ip
}

resource "panos_arp" "x" {
    template = panos_panorama_ethernet_interface.x.template
    interface_type = "ethernet"
    interface_name = panos_panorama_ethernet_interface.x.name
    ip = "10.5.6.7"
    mac_address = "00:30:48:52:11:22"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "x" {
    name = "template one"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"

    lifecycle {
        create_before_destroy = true
    }
}


# Panorama aggregate interface example.
data "panos_arp" "example2" {
    template = panos_arp.y.template
    interface_type = panos_arp.y.interface_type
    interface_name = panos_arp.y.interface_name
    ip = panos_arp.y.ip
}

resource "panos_arp" "y" {
    template = panos_panorama_aggregate_interface.y.template
    interface_type = "aggregate-ethernet"
    interface_name = panos_panorama_aggregate_interface.y.name
    ip = "10.5.6.7"
    mac_address = "00:30:48:52:22:33"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "y" {
    name = "template two"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_aggregate_interface" "y" {
    template = panos_panorama_template.y.name
    name = "ae1"
    vsys = "vsys1"
    mode = "layer3"

    lifecycle {
        create_before_destroy = true
    }
}


# Panorama VLAN interface example.
# Since all VLAN interfaces are subinterfaces and not top level interfaces,
# the `interface_name` param should be left empty.
data "panos_arp" "example3" {
    template = panos_arp.z.template
    interface_type = panos_arp.z.interface_type
    subinterface_name = panos_arp.z.subinterface_name
    ip = panos_arp.z.ip
}

resource "panos_arp" "z" {
    template = panos_panorama_vlan_interface.z.template
    interface_type = "vlan"
    subinterface_name = panos_panorama_vlan_interface.z.name
    ip = "10.5.6.7"
    mac_address = "00:30:48:52:33:44"

    lifecycle {
        create_before_destroy = true
    }
}

resourcee "panos_panorama_template" "z" {
    name = "template three"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_vlan_interface" "z" {
    template = panos_panorama_template.z.name
    name = "vlan.42"
    vsys = "vsys1"

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `interface_type` - The interface type.  Valid values are `ethernet` (default),
  `aggregate-ethernet`, or `vlan`.
* `interface_name` - The interface name (leave this empty for VLAN interfaces).
* `subinterface_name` - The subinterface name.
* `ip` - (Required) The IP address.

## Attribute Reference

The following attributes are supported:

* `mac_address` - The MAC address.
* `interface` - (`interface_type`=`vlan`) The interface.
