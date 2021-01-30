---
page_title: "panos: panos_device_group_parent"
subcategory: "Panorama"
---


# panos_device_group_parent

Manage the device group parent for a given device group.

**NOTE:** This is for Panorama only.


## Example Usage

```hcl
# Example 1:  Move "Group B" under "Group A".
resource "panos_device_group_parent" "example1" {
    device_group = panos_panorama_template.b.name
    parent = panos_panorama_template.a.name
}

resource "panos_panorama_template" "a" {
    name = "Group A"
}

resource "panos_panorama_template" "b" {
    name = "Group B"
}


# Example 2:  Ensure that "Group C" is under "shared" and has no parent.
resource "panos_device_group_parent" "example2" {
    device_group = panos_panorama_template.c.name
}

resource "panos_panorama_template" "c" {
    name = "Group C"
}
```


## Argument Reference

The following arguments are supported:

* `device_group` - (Required) The device group whose parent you intent to set.
* `parent` - The parent device group name.  Leaving this empty / unspecified
  means to move this device group under the "shared" device group.
