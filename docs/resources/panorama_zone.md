---
page_title: "panos: panos_panorama_zone"
subcategory: "Panorama Networking"
---

# panos_panorama_zone

This resource allows you to add/update/delete zones on Panorama for both
templates and template stacks.

This resource has some overlap with the `panos_panorama_zone_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_zone` spec does not define the
`interfaces` field.


## Import Name

```
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_zone" "example" {
    name = "myZone"
    template = panos_panorama_template.tmpl1.name
    mode = "layer3"
    interfaces = [
        panos_panorama_ethernet_interface.e2.name,
        panos_panorama_ethernet_interface.e3.name,
    ]
    enable_user_id = true
    exclude_acls = ["192.168.0.0/16"]
}

resource "panos_panorama_template" "tmpl1" {
    name = "MyTemplate"
}

resource "panos_panorama_ethernet_interface" "e2" {
    template = panos_panorama_template.tmpl1.name
    name = "ethernet1/2"
    mode = "layer3"
}

resource "panos_panorama_ethernet_interface" "e3" {
    template = panos_panorama_template.tmpl1.name
    name = "ethernet1/3"
    mode = "layer3"
}
```

## Argument Reference

At least one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The zone's name.
* `vsys` - (Optional) The vsys to put the zone into (default: `vsys1`).
* `mode` - (Required) The zone's mode.  This can be `layer3`, `layer2`,
  `virtual-wire`, `tap`, or `tunnel`.
* `zone_profile` - (Optional) The zone protection profile.
* `log_setting` - (Optional) Log setting.
* `enable_user_id` - (Optional) Boolean to enable user identification.
* `interfaces` - (Optional) List of interfaces to associated with this zone.  If
  you are going to use the `panos_panorama_zone_entry` resource then this param
  should be left unspecified.
* `include_acls` - (Optional) Users from these addresses/subnets will
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
* `exclude_acls` - (Optional) Users from these addresses/subnets will not
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
