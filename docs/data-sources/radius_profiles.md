---
page_title: "panos: panos_radius_profiles"
subcategory: "Device"
---

# panos_radius_profiles

Gets the list of Radius profiles.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_radius_profiles" "example" {}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.

NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
