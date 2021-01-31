---
page_title: "panos: panos_ospf_area_interfaces"
subcategory: "Network"
---

# panos_ospf_area_interfaces

Gets a list of OSPF area interfaces.


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_area_interfaces" "example" {
    template = "my template"
    virtual_router = "my virtual router"
    ospf_area = "10.2.3.1"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.
* `ospf_area` - (Required) OSPF area name.


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
