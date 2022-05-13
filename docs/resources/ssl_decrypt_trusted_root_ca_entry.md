---
page_title: "panos: panos_ssl_decrypt_trusted_root_ca_entry"
subcategory: "Device"
---

# panos_ssl_decrypt_trusted_root_ca_entry

This resource manages the SSL decrypt settings.

This resource has some overlap with the
[`panos_ssl`](ssl_decrypt.html)
resource.  If you want to use this resource with the other one, then make sure that
your `trusted_root_cas` param is left undefined.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<template>:<template_stack>:<vsys>:<certificate_name>
```


## Example Usage

```hcl
resource "panos_ssl_decrypt_trusted_root_ca_entry" "example" {
    vsys = panos_ssl_decrypt.x.vsys
    template = panos_ssl_decrypt.x.template
    template_stack = panos_ssl_decrypt.x.template_stack
    certificate_name = panos_certificate_import.trust.name
}

resource "panos_ssl_decrypt" "x" {}

resource "panos_certificate_import" "trust" {
    name = "tfTrust"
    pem {
        certificate = file("cert.pem")
        private_key = file("key.pem")
        passphrase = "secret"
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

* `certificate_name` - (Required) The certificate name.
