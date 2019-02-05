---
layout: "panos"
page_title: "panos: panos_panorama_zone_entry"
sidebar_current: "docs-panos-panorama-resource-zone-entry"
description: |-
  Manages a specific interface in a Panorama zone.
---

# panos_panorama_zone_entry

This resource allows you to add/update/delete a specific interface in a Panorama
zone.

This resource has some overlap with the `panos_panorama_zone`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_zone` spec does not define the
`interfaces` field.

This is the appropriate resource to use if you have a pre-existing zone
in Panorama and don't want Terraform to delete it on `terraform destroy`.


## Import Name

```
<template>:<template_stack>:<vsys>:<zone>:<mode>:<interface>
```


## Example Usage

```hcl
resource "panos_panorama_template" "t" {
    name = "myTemplate"
}

resource "panos_panorama_ethernet_interface" "e5" {
    template = "${panos_panorama_template.t.name}"
    name = "ethernet1/5"
    mode = "layer3"
}

resource "panos_panorama_zone" "z" {
    template = "${panos_panorama_template.t.name}"
    name = "exZone"
    mode = "layer3"
}

resource "panos_panorama_zone_entry" "example" {
    template = "${panos_panorama_template.t.name}"
    zone = "${panos_panorama_zone.z.name}"
    mode = "${panos_panorama_zone.z.mode}"
    interface = "${panos_panorama_ethernet_interface.e5.name}"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
* `vsys` - (Optional) The vsys (default: `vsys1`).
* `zone` - (Required) The zone's name.
* `mode` - (Optional) The mode.  Can be `layer3` (default), `layer2`,
  `virtual-wire`, `tap`, or `external`.
* `interface` - (Required) The interface's name.
