---
page_title: "panos: panos_ospf_exports"
subcategory: "Network"
---

# panos_ospf_exports

Gets the list of OSPF export rules.


## Example Usage

```hcl
# Panorama example.
data "panos_ospf_exports" "example" {
    template = "my template"
    virtual_router = "my virtual router"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
