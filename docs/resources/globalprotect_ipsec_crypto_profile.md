---
page_title: "panos: panos_globalprotect_ipsec_crypto_profile"
subcategory: "Network"
---

# panos_globalprotect_ipsec_crypto_profile

Manages a GlobalProtect IPSec crypto profile.


## PAN-OS

NGFW and Panorama.

## Import Name

```
<template>:<template_stack>:<name>
```

If `template` and `template_stack` are empty, use:
```
::<name>
```

## Example Usage

```hcl
resource "panos_globalprotect_ipsec_crypto_profile" "example" {
    name = "fromTerraform"
    encryptions = [
        "aes-128-gcm",
        "aes-256-gcm",
    ]
    authentications = ["sha1"]

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


The following arguments are supported:

* `name` - (Required) The name.
* `encryptions - (List of string) The encryptions.  Valid values are `"aes-128-cbc"`, `"aes-128-gcm"`, and `"aes-256-gcm"`.
* `authentications` - (List of string) The authentication algorithms.  Valid values are `"sha1"`.
