---
page_title: "panos: panos_ospf_area_virtual_links"
subcategory: "Network"
---

# panos_ospf_area_virtual_links

Gets a list of OSPF area virtual links.


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_area_virtual_links" "example" {
    template = "my template"
    virtual_router = "the virtual router name"
    ospf_area = "10.10.10.10"
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
