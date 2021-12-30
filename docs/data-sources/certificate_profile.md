---
page_title: "panos: panos_certificate_profile"
subcategory: "Device"
---

# panos_certificate_profile

Gets information on a certificate profiles.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_certificate_profile" "example" {
    name = "manual"
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

## Attribute Reference

The following attributes are supported:

* `username_field` - Username field.  Valid values are an empty string for `None`,
  `subject`, or `subject-alt`.
* `username_field_value` - The value.
* `domain` - User domain.
* `use_crl` - (bool) Use CRL.
* `use_ocsp` - (bool) Use OCSP.
* `crl_receive_timeout` - (int) CRL receive timeout in seconds.
* `ocsp_receive_timeout` - (int) OCSP receive timeout in seconds.
* `certificate_status_timeout` - (int) Certificate status timeout in seconds.
* `block_unknown_certificate` - (bool) Block session if certificate status
  is unknown.
* `block_certificate_timeout` - (bool) Block session if certificate status
  cannot be retrieved within timeout.
* `block_unauthenticated_certificate` - (bool) Block session if the certificate
  was not issued to the authenticating device.
* `block_expired_certificate` - (bool) Block sessions with expired certificates.
* `ocsp_exclude_nonce` - (bool) OCSP exclude nonce.
* `certificate` - (repeated) List of CA certificates, defined below.

`certificate` supports the following attributes:

* `name` - The name.
* `default_ocsp_url` - Default OCSP URL.
* `ocsp_verify_certificate` - OCSP verify certificate.
* `template_name` - Template name/OID.
