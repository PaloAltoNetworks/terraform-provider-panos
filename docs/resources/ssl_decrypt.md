---
page_title: "panos: panos_ssl_decrypt"
subcategory: "Device"
---

# panos_ssl_decrypt

This resource manages the SSL decrypt settings.

This resource has some overlap with the
[`panos_ssl_decrypt_trusted_root_ca_entry`](ssl_decrypt_trusted_root_ca_entry.html)
resource.  If you want to use this resource with the other one, then make sure that
your `trusted_root_cas` param is left undefined.


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
resource "panos_ssl_decrypt" "example" {
    forward_trust_certificate_rsa = panos_certificate_import.trust.name
    forward_untrust_certificate_rsa = panos_certificate_import.untrust.name
}

resource "panos_certificate_import" "trust" {
    name = "tfTrust"
    pem {
        certificate = file("cert.pem")
        private_key = file("key.pem")
        passphrase = "secret"
    }
}

resource "panos_certificate_import" "untrust" {
    name = "tfUntrust"
    pem {
        certificate = file("untrust-cert.pem")
        private_key = file("untrust-key.pem")
        passphrase = "foobar"
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

* `forward_trust_certificate_rsa` - Forward trust RSA certificate.
* `forward_trust_certificate_ecdsa` - Forward trust ECDSA certificate.
* `forward_untrust_certificate_rsa` - Forward untrust RSA certificate.
* `forward_untrust_certificate_ecdsa` - Forward untrust ECDSA certificate.
* `root_ca_excludes` - List of root CA excludes.
* `trusted_root_cas` - List of trusted root CAs.
* `disabled_predefined_exclude_certificates` - List of disabled predefined
  exclude certificates.
* `ssl_decrypt_exclude_certificate` - (repeatable) List of SSL decrypt exclude
  certificates specs (specified below).


`ssl_decrypt_exclude_certificate` sections support the following arguments:

* `name` - (Required) The name.
* `description` - The description.
* `exclude` - (bool) Exclude or not.
