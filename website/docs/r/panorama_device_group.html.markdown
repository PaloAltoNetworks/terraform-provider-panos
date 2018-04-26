---
layout: "panos"
page_title: "panos: panos_panorama_device_group"
sidebar_current: "docs-panos-panorama-resource-device-group"
description: |-
  Manages Panorama device groups.
---

# panos_panorama_device_group

This resource allows you to add/update/delete Panorama device groups.

This resource has some overlap with the `panos_panorama_device_group_entry`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_device_group` spec does not define any
`device` blocks, and just stays as "computed".

This is the appropriate resource to use if `terraform destroy` should delete
the device group.

## Example Usage

```hcl
resource "panos_panorama_device_group" "example" {
    name = "my device group"
    description = "description here"
    device {
        serial = "00112233"
    }
    device {
        serial = "44556677"
        vsys_list = ["vsys1", "vsys2"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The device group's name.
* `description` - (Optional) The device group's description.
* `device` - The device definition (see below).

The following arguments are valid for each `device` section:

* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.
