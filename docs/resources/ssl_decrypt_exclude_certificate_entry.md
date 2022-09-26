---
page_title: "panos: panos_ssl_decrypt_exclude_certificate_entry"
subcategory: "Device"
---

# panos_ssl_decrypt_exclude_certificate_entry

This resource manages the exclude certificates within the SSL decrypt settings.

This resource has some overlap with the [`panos_ssl`](ssl_decrypt.html) resource.  If you want to use this resource with the other one, then make sure that no `ssl_decrypt_exclude_certificate` sections are defined.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_ssl_decrypt_exclude_certificate_entry" "example" {
    vsys = panos_ssl_decrypt.x.vsys
    template = panos_ssl_decrypt.x.template
    template_stack = panos_ssl_decrypt.x.template_stack
    name = "*.example.com"
    description = "example"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_ssl_decrypt" "x" {}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).


The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `exclude` - (bool) Exclude or not (default: `true`).
