---
page_title: "panos: panos_certificate_profiles"
subcategory: "Device"
---

# panos_certificate_profiles

Gets the list of certificate profiles.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_certificate_profile" "example" {}
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
