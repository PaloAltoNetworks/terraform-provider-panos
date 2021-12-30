---
page_title: "panos: panos_panorama_ipsec_crypto_profile"
subcategory: "Network"
---

# panos_panorama_ipsec_crypto_profile

This resource allows you to add/update/delete Panorama IPSec crypto profiles
for both templates and template stacks.


## PAN-OS

Panorama


## Import Name

```
<template>:<template_stack>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_ipsec_crypto_profile" "example" {
    template = panos_panorama_template.t.name
    name = "example"
    authentications = ["md5", "sha384"]
    encryptions = ["des", "aes-128-cbc"]
    dh_group = "group14"
    lifetime_type = "hours"
    lifetime_value = 4
    lifesize_type = "mb"
    lifesize_value = 1
}

resource "panos_panorama_template" "t" {
    name = "my template"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The object's name
* `protocol` - (Optional) The protocol.  Valid values are `esp` (the default)
  or `ah`
* `authentications` - (Required, list) - List of authentication types.
* `encryptions` - (Required, list) - List of encryption types.  Valid values
  are `des`, `3des`, `aes-128-cbc`, `aes-192-cbc`, `aes-256-cbc`, `aes-128-gcm`,
  `aes-256-gcm`, and `null`.  Note that the "gcm" values are only available in
  PAN-OS 7.0+.
* `dh_group` - (Optional) The DH group value.  Valid values should start with
  the string `group`.
* `lifetime_type` - (Optional) The lifetime type.  Valid values are `seconds`,
  `minutes`, `hours` (the default), or `days`.
* `lifetime_value` - (Optional, int) The lifetime value.
* `lifesize_type` - (Optional) The lifesize type.  Valid values are `kb`, `mb`,
  `gb`, or `tb`.
* `lifesize_value` - (Optional, int) the lifesize value.
