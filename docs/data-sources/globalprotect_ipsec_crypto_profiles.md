---
page_title: "panos: panos_globalprotect_ipsec_crypto_profiles"
subcategory: "Network"
---

# panos_globalprotect_ipsec_crypto_profiles

Gets information on a GlobalProtect IPSec crypto profile.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_globalprotect_ipsec_crypto_profiles" "example" {}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.

The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `encryptions - (List of string) The encryptions.
* `authentications` - (List of string) The authentication algorithms.
