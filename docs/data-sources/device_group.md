---
page_title: "panos: panos_device_group"
subcategory: "Panorama"
---

# panos_device_group

Retrieve information on the specified device group.


## PAN-OS

Panorama only.


## Example Usage

```hcl
data "panos_device_group" "example" {
    name = "my device group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The device group's name.


## Attribute Reference

The following attributes are supported:

* `description` - (Optional) The device group's description.
* `device` - The device definition (see below).

The following arguments are valid for each `device` section:

* `serial` - The serial number of the firewall.
* `vsys_list` - A subset of all available vsys on the firewall
  that should be in this device group.
