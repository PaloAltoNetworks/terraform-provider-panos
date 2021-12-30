---
page_title: "panos: panos_ssl_decrypt"
subcategory: "Device"
---

# panos_ssl_decrypt

Retrieves information on the SSL decrypt settings.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_ssl_decrypt" "example" {}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW / Panorama:

* `vsys` - The vsys (default: `shared`).


## Attribute Reference

The following attributes are supported:

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


`ssl_decrypt_exclude_certificate` sections support the following attributes:

* `name` - The name.
* `description` - The description.
* `exclude` - (bool) Exclude or not.
