---
page_title: "panos: panos_ssl_tls_service_profile"
subcategory: "Device"
---

# panos_ssl_tls_service_profile

Gets information on a SSL TLS service profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_ssl_tls_service_profiles" "example" {
    name = "fromTerraform"
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

* `certificate` - SSL certificate file name.
* `min_version` - Minimum TLS protocol version.
* `max_version` - Maximum TLS protocol version.
* `allow_algorithm_rsa` - (bool) Allow algorithm RSA.
* `allow_algorithm_dhe` - (bool) Allow algorithm DHE.
* `allow_algorithm_ecdhe` - (bool) Allow algorithm ECDHE.
* `allow_algorithm_3des` - (bool) Allow algorithm 3DES.
* `allow_algorithm_rc4` - (bool) Allow algorithm RC4.
* `allow_algorithm_aes_128_cbc` - (bool) Allow algorithm AES-128-CBC.
* `allow_algorithm_aes_256_cbc` - (bool) Allow algorithm AES-256-CBC.
* `allow_algorithm_aes_128_gcm` - (bool) Allow algorithm AES-128-GCM.
* `allow_algorithm_aes_256_gcm` - (bool) Allow algorithm AES-256-GCM.
* `allow_authentication_sha1` - (bool) Allow authentication SHA1.
* `allow_authentication_sha256` - (bool) Allow authentication SHA256.
* `allow_authentication_sha384` - (bool) Allow authentication SHA384.
