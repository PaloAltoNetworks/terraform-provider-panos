---
page_title: "panos: panos_arps"
subcategory: "Network"
---

# panos_arps

Gets the list of ARP configs attached to an interface.


## Example Usage

```hcl
# Panorama ethernet interface example.
data "panos_arps" "example" {
    template = "some template"
    interface_type = "ethernet"
    interface_name = "ethernet1/1"
    subinterface_name = "ethernet1/1.42"
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


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
