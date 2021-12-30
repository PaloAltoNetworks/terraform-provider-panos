---
page_title: "panos: panos_device_group_parent"
subcategory: "Panorama"
---


# panos_device_group_parent

Manage the device group parent for a given device group.


## PAN-OS

Panorama.


## Import Name

```
<device_group>
```


## Example Usage

```hcl
# Example 1:  Move "Group B" under "Group A".
resource "panos_device_group_parent" "example1" {
    device_group = panos_device_group.b.name
    parent = panos_device_group.a.name
}

resource "panos_device_group" "a" {
    name = "Group A"
}

resource "panos_device_group" "b" {
    name = "Group B"
}
```

```hcl
# Example 2:  Ensure that "Group C" is under "shared" and has no parent.
resource "panos_device_group_parent" "example2" {
    device_group = panos_device_group.c.name
}

resource "panos_device_group" "c" {
    name = "Group C"
}
```


## Argument Reference

The following arguments are supported:

* `device_group` - (Required) The device group whose parent you intent to set.
* `parent` - The parent device group name.  Leaving this empty / unspecified
  means to move this device group under the "shared" device group.
