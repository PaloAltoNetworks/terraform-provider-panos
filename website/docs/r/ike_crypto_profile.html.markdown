---
layout: "panos"
page_title: "panos: panos_ike_crypto_profile"
sidebar_current: "docs-panos-resource-ike-crypto-profile"
description: |-
  Manages IKE crypto profiles.
---

# panos_ike_crypto_profile

This resource allows you to add/update/delete IKE crypto profiles.

## Example Usage

```hcl
resource "panos_ike_crypto_profile" "example" {
    name = "example"
    dh_groups = ["group1", "group2"]
    authentications = ["md5", "sha1"]
    encryptions = ["des"]
    lifetime_value = 8
    authentication_multiple = 3
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name
* `dh_groups` - (Required, list) List of DH Group entries.  Values should
  have a prefix if `group`.
* `authentications` - (Required, list) List of authentication types.  This c
* `encryptions` - (Required, list) List of encryption types.  Valid values
  are `des`, `3des`, `aes-128-cbc`, `aes-192-cbc`, and `aes-256-cbc`.
* `lifetime_type` - (Optional) The lifetime type.  Valid values are `seconds`,
  `minutes`, `hours` (the default), and `days`.
* `lifetime_value` - (Optional, int) The lifetime value.
* `authentication_multiple` - (Optional, PAN-OS 7.0+, int) IKEv2 SA
  reauthentication interval equals authetication-multiple * rekey-lifetime; 0
  means reauthentication is disabled.
