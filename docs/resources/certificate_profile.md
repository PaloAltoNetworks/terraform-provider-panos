---
page_title: "panos: panos_certificate_profile"
subcategory: "Device"
---

# panos_certificate_profile

This resource allows you to add/update/delete certificate profiles.


## PAN-OS

NGFW and Panorama.


## Import Name

```
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_certificate_profile" "x" {
    name = "manual"
    username_field = "subject"
    username_field_value = "common-name"
    domain = "blah"
    certificate {
        name = "myCert"
        default_ocsp_url = "https://hello.example.com/default"
        ocsp_verify_certificate = ""
        template_name = ""
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

* `name` - (Required) The name.
* `username_field` - Username field.  Valid values are an empty string for `None`
  (default), `subject`, or `subject-alt`.
* `username_field_value` - The value.  Common settings are `common-name` for
  `username_field="subject"`, or `email` or `principal-name` for
  `username_field="subject-alt"`.
* `domain` - User domain.
* `use_crl` - (bool) Use CRL.
* `use_ocsp` - (bool) Use OCSP.
* `crl_receive_timeout` - (int) CRL receive timeout in seconds (default: `5`).
* `ocsp_receive_timeout` - (int) OCSP receive timeout in seconds (default: `5`).
* `certificate_status_timeout` - (int) Certificate status timeout in
  seconds (default: `5`).
* `block_unknown_certificate` - (bool) Block session if certificate status
  is unknown.
* `block_certificate_timeout` - (bool) Block session if certificate status
  cannot be retrieved within timeout.
* `block_unauthenticated_certificate` - (bool) Block session if the certificate
  was not issued to the authenticating device.
* `block_expired_certificate` - (bool) Block sessions with expired certificates.
* `ocsp_exclude_nonce` - (bool) OCSP exclude nonce.
* `certificate` - (repeated) List of CA certificates, defined below.

`certificate` supports the following arguments:

* `name` - (Required) The name.
* `default_ocsp_url` - Default OCSP URL.
* `ocsp_verify_certificate` - OCSP verify certificate.
* `template_name` - Template name/OID.
