---
page_title: "panos: panos_device_groups"
subcategory: "Panorama"
---

# panos_device_groups

Gets the list of device groups.


## Example Usage

```hcl
data "panos_device_groups" "example" {}
```


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
