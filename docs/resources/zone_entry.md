---
page_title: "panos: panos_zone_entry"
subcategory: "Network"
---

# panos_zone_entry

This resource allows you to add/update/delete a specific interface in a zone.

This resource has some overlap with the [`panos_zone`](zone.html)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_zone` spec does not define the
`interfaces` field.

This is the appropriate resource to use if you have a pre-existing zone
and don't want Terraform to delete it on `terraform destroy`.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_zone_entry`


## Import Name

```
<template>:<template_stack>:<vsys>:<zone>:<mode>:<interface>
```


## Example Usage

```hcl
resource "panos_zone_entry" "example" {
    zone = panos_zone.z.name
    mode = panos_zone.z.mode
    interface = panos_ethernet_interface.e5.name
}

resource "panos_ethernet_interface" "e5" {
    name = "ethernet1/5"
    mode = "layer3"
}

resource "panos_zone" "z" {
    name = "exZone"
    mode = "layer3"
}
```

## Argument Reference

Panorama specific arguments (one of these must be specified):

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `vsys` - (Optional) The vsys (default: `vsys1`).
* `zone` - (Required) The zone's name.
* `mode` - (Optional) The mode.  Can be `layer3` (default), `layer2`,
  `virtual-wire`, `tap`, or `external`.
* `interface` - (Required) The interface's name.
