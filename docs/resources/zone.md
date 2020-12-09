---
page_title: "panos: panos_zone"
subcategory: "Firewall Networking"
---

# panos_zone

This resource allows you to add/update/delete zones.

This resource has some overlap with the `panos_zone_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_zone` spec does not define the
`interfaces` field.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_zone" "example" {
    name = "myZone"
    mode = "layer3"
    interfaces = [
        panos_ethernet_interface.e1.name,
        panos_ethernet_interface.e5.name,
    ]
    enable_user_id = true
    exclude_acls = ["192.168.0.0/16"]
}

resource "panos_ethernet_interface" "e1" {
    name = "ethernet1/1"
    mode = "layer3"
}

resource "panos_ethernet_interface" "e5" {
    name = "ethernet1/5"
    mode = "layer3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The zone's name.
* `vsys` - (Optional) The vsys to put the zone into (default: `vsys1`).
* `mode` - (Required) The zone's mode.  This can be `layer3`, `layer2`,
  `virtual-wire`, `tap`, or `tunnel`.
* `zone_profile` - (Optional) The zone protection profile.
* `log_setting` - (Optional) Log setting.
* `enable_user_id` - (Optional) Boolean to enable user identification.
* `interfaces` - (Optional) List of interfaces to associated with this zone.  Leave
  this undefined if you want to use `panos_zone_entry` resources.
* `include_acls` - (Optional) Users from these addresses/subnets will
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
* `exclude_acls` - (Optional) Users from these addresses/subnets will not
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
