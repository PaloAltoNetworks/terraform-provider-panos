---
layout: "panos"
page_title: "panos: panos_panorama_device_group_entry"
sidebar_current: "docs-panos-panorama-resource-device-group-entry"
description: |-
  Manages a specific device in a Panorama device group.
---

# panos_panorama_device_group_entry

This resource allows you to add/update/delete a specific device in a Panorama
device group.

This resource has some overlap with the `panos_panorama_device_group`
resource.  If you want to use this resource with the other one, then make
sure that your `panos_panorama_device_group` spec does not define any
`device` blocks, and just stays as "computed".

This is the appropriate resource to use if you have a pre-existing device group
in Panorama and don't want Terraform to delete it on `terraform destroy`.

An interesting side effect of the underlying XML API - if the device group does
not already exist, then this resource can actually create it.  However, since
only the single entry for the specific serial number is deleted, then a
`terraform destroy` would not remove the device group itself in this situation.


## Import Name

```
<device_group>:<serial>
```


## Example Usage

```hcl
# Example for a virtual firewall.
resource "panos_panorama_device_group_entry" "example1" {
    device_group = "my device group"
    serial = "00112233"
}

# Example for a physical firewall with multi-vsys enabled.
resource "panos_panorama_device_group_entry" "example2" {
    device_group = "my device group"
    serial = "44556677"
    vsys_list = ["vsys1", "vsys2"]
}
```

## Argument Reference

The following arguments are supported:

* `device_group` - (Required) The device group's name.
* `serial` - (Required) The serial number of the firewall.
* `vsys_list` - (Optional) A subset of all available vsys on the firewall
  that should be in this device group.  If the firewall is a virtual firewall,
  then this parameter should just be omitted.
