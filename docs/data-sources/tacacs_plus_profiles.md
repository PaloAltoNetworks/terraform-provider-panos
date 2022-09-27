---
page_title: "panos: panos_tacacs_plus_profiles"
subcategory: "Device"
---

# panos_tacacs_plus_profiles

Gets the list of TACACS+ profiles.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_tacacs_plus_profiles" "example" {}
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
