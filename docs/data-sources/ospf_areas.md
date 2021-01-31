---
page_title: "panos: panos_ospf_areas"
subcategory: "Network"
---

# panos_ospf_areas

Gets a list of OSPF areas.


## Example Usage

```hcl
data "panos_ospf_areas" "example" {
    template = "my template"
    virtual_router = "some virtual router"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.


## Attribute Reference

The following attributes are available:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
