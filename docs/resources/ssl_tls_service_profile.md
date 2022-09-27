---
page_title: "panos: panos_ssl_tls_service_profile"
subcategory: "Device"
---

# panos_ssl_tls_service_profile

Manages a SSL TLS service profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_ssl_tls_service_profiles" "example" {
    name = "fromTerraform"
    certificate = "myCert"
    min_version = "tls1-1"
    allow_algorithm_ecdhe = false
    allow_algorithm_aes_128_gcm = false
    allow_algorithm_aes_256_gcm = false

    lifecycle {
        create_before_destroy = true
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
* `certificate` - (Required) SSL certificate file name.
* `min_version` - Minimum TLS protocol version.  Valid values are `"tls1-0"` (default), `"tls1-1"`, and `"tls1-2"`.
* `max_version` - Maximum TLS protocol version.  Valid values are `"tls1-0"`, `"tls1-1"`, `"tls1-2"`, and `max` (default).
* `allow_algorithm_rsa` - (bool) Allow algorithm RSA (default: `true`).
* `allow_algorithm_dhe` - (bool) Allow algorithm DHE (defualt: `true`).
* `allow_algorithm_ecdhe` - (bool) Allow algorithm ECDHE (default: `true`).
* `allow_algorithm_3des` - (bool) Allow algorithm 3DES (default: `true`).
* `allow_algorithm_rc4` - (bool) Allow algorithm RC4 (default: `true`).
* `allow_algorithm_aes_128_cbc` - (bool) Allow algorithm AES-128-CBC (default: `true`).
* `allow_algorithm_aes_256_cbc` - (bool) Allow algorithm AES-256-CBC (defualt: `true`).
* `allow_algorithm_aes_128_gcm` - (bool) Allow algorithm AES-128-GCM (default: `true`).
* `allow_algorithm_aes_256_gcm` - (bool) Allow algorithm AES-256-GCM (default: `true`).
* `allow_authentication_sha1` - (bool) Allow authentication SHA1 (default: `true`).
* `allow_authentication_sha256` - (bool) Allow authentication SHA256 (default: `true`).
* `allow_authentication_sha384` - (bool) Allow authentication SHA384 (default: `true`).
