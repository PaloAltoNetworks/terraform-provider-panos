---
page_title: "panos: panos_device_group_parent"
subcategory: "Panorama"
---


# panos_device_group_parent

Retrieves the device group hierarchy as a map.

**NOTE:** This is for Panorama only.


## Example Usage

```hcl
data "panos_device_group_parent" "example" {}
```


## Attribute Reference

The following attributes are supported:

* `total` - (int) Total number of entries (device groups).
* `entries` - (map of strings) Map of strings where the key is the device
  group name and the value is the parent for that device group.  An empty
  string for the value means that the parent is the "shared" device group.
