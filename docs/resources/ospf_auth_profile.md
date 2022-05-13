---
page_title: "panos: panos_ospf_auth_profile"
subcategory: "Network"
---

# panos_ospf_auth_profile

Manages OSPF auth profile config attached to a virtual router.


## Import Name

NGFW:

```shell
<virtual_router>:<name>
```

Panorama:

```shell
<template>::<virtual_router>:<name>
```


## Example Usage

```hcl
# Panorama example.
resource "panos_ospf_auth_profile" "example" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = "my profile"
    password = "secret"
}

resource "panos_ospf" "x" {
    template = panos_panorama_template.x.name
    virtual_router = panos_panorama_virtual_router.x.name
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = "my virtual router"
}       

resource "panos_panorama_template" "x" {
    name = "my template"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.
* `name` - (Required) The export rule name.
* `auth_type` - The auth type.  Valid values are `password` (default) or `md5`.
* `password` - The simple password.
* `md5_key` - (list) List of md5_key specs, as defined below.

`md5_key` supports the following arguments:

* `key_id` - (Required, int) MD5 key ID.
* `key` - (Required) MD5 key.
* `preferred` - (bool) Preferred key.


## Attribute Reference

The following attributes are supported:

* `password_enc` - Encrypted simple password.
* `md5_keys_enc` - (list) List of encrypted/unencrypted MD5 keys.
