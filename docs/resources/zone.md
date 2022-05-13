---
page_title: "panos: panos_zone"
subcategory: "Network"
---

# panos_zone

This resource allows you to add/update/delete zones.

This resource has some overlap with the [`panos_zone_entry`](zone_entry.html)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_zone` spec does not define the
`interfaces` field.


## PAN-OS

NGFW and Panorama


## Aliases

* `panos_panorama_zone`


## Import Names

```shell
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

### Firewall Example

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

### Panorama Example

```hcl
resource "panos_zone" "example" {
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

Panorama specific arguments (one of these must be specified):

* `template` - The template name.
* `template_stack` - The template stack name.

NGFW / Panorama:

* `vsys` - The vsys (default: `vsys1`).

The following arguments are supported:

* `name` - (Required) The zone's name.
* `mode` - (Required) The zone's mode.  This can be `layer3`, `layer2`,
  `virtual-wire`, `tap`, or `tunnel`.
* `zone_profile` - The zone protection profile.
* `log_setting` - Log setting.
* `enable_user_id` - Boolean to enable user identification.
* `interfaces` - List of interfaces to associated with this zone.  Leave
  this undefined if you want to use [`panos_zone_entry`](zone_entry.html) resources.
* `include_acls` - Users from these addresses/subnets will
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
* `exclude_acls` - Users from these addresses/subnets will not
  be identified.  This can be an address object, an address group, a single
  IP address, or an IP address subnet.
