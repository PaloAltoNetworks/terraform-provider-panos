---
page_title: "panos: panos_ssl_decrypt_exclude_certificate_entry"
subcategory: "Device"
---

# panos_ssl_decrypt_exclude_certificate_entry

This resource manages the SSL Decrypt Exclusion setting.

This resource has some overlap with the
[`panos_ssl`](ssl_decrypt.html)
resource.  If you want to use this resource with the other one, then make sure that
your `ssl_decrypt_exclude_certificate` param is left undefined.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<template>:<template_stack>:<vsys>
```


## Example Usage

```hcl
resource "panos_ssl_decrypt_exclude_certificate_entry" "example" {
    ssl_decrypt_exclude_certificate {
      name = "*.example.com"
      description = "example"
      exclude = "true"
    }
}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).


The following arguments are supported:

* `ssl_decrypt_exclude_certificate` - (repeatable) List of SSL decrypt exclude
  certificates specs (specified below).


`ssl_decrypt_exclude_certificate` sections support the following arguments:

* `name` - (Required) The name.
* `description` - The description.
* `exclude` - (bool) Exclude or not.
