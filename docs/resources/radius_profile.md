---
page_title: "panos: panos_radius_profile"
subcategory: "Device"
---

# panos_radius_profile

Manages a Radius profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_radius_profiles" "example" {
    name = "fromTerraform"
    timeout = 4
    retries = 5
    protocol {
        eap_ttls_with_pap {
            make_outer_identity_anonymous = true
            certificate_profile = "someCertProfile"
        }
    }
    server {
        name = "first"
        server = "first.example.com"
        secret = "secret"
    }
    server {
        name = "second"
        server = "192.168.0.5"
        secret = "anotherSecret"
        port = 1234
    }

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
* `admin_use_only` - (bool) Administrator use only.
* `timeout` - (int) Timeout in seconds (default: `3`).
* `retries` - (int) Number of attempts before giving up authentication (default: `3`).
* `protocol` - (Required, PAN-OS 8.0+) Authentication protocol settings spec, as defined below.
* `server` - List of server specs, as defined below.

`protocol` supports the following arguments:

* `chap` - (bool) CHAP.
* `pap` - (bool) PAP.
* `auth` - (bool, PAN-OS 8.0 only) Auto.
* `peap_mschap_v2` - PEAP-MSCHAPv2 spec, as defined below.
* `peap_with_gtc` - PEAP with GTC spec, as defined below.
* `eap_ttls_with_pap` - EAP-TTLS with PAP spec, as defined below.

`protocol.peap_mschap_v2` supports the following arguments:

* `make_outer_identity_anonymous` - (bool) Make outer identity anonymous.
* `allow_expired_password_change` - (bool) Allow users to change passwords after expiry.
* `certificate_profile` - (Required) Certificate profile for verifying the Radius server.

`protocol.peap_with_gtc` supports the following arguments:

* `make_outer_identity_anonymous` - (bool) Make outer identity anonymous.
* `certificate_profile` - (Required) Certificate profile for verifying the Radius server.

`protocol.eap_ttls_with_pap` supports the following arguments:

* `make_outer_identity_anonymous` - (bool) Make outer identity anonymous.
* `certificate_profile` - (Required) Certificate profile for verifying the Radius server.

`server` supports the following arguments:

* `name` - (Required) The name.
* `server` - (Required) Server hostname or IP address.
* `secret` - (Required) Shared secret for Radius communication.
* `port` - (int) Radius server port number (default: `1812`).
