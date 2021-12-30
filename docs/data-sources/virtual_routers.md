---
page_title: "panos: panos_virtual_routers"
subcategory: "Network"
---

# panos_virtual_routers

Retrieve the list of virtual routers.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_virtual_router" "example" {}
```


## Argument Reference

Panorama (currently only templates can have virtual routers):

* `template` - The template.
* `template_stack` - The template stack.


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
