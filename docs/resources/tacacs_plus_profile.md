---
page_title: "panos: panos_tacacs_plus_profile"
subcategory: "Device"
---

# panos_tacacs_plus_profile

Manages a TACACS+ profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_tacacs_plus_profiles" "example" {
    name = "fromTerraform"
    admin_use_only = true
    timeout = 5
    use_single_connection = true
    protocol {
        pap = true
    }
    server {
        name = "second"
        server = "tacacs_plus.example.com
        secret = "mySecret"
    }
    server {
        name = "second"
        server = "192.168.1.50"
        secret = "drowssap"
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
* `use_single_connection` - (bool) Use single connectino for all authentication."
* `protocol` - Protocol spec, as defined below.
* `server` - List of TACACS+ server specs, as defined below.

`protocol` supports the following arguments:

* `chap` - (bool) CHAP.
* `pap` - (bool) PAP.
* `auto` - (bool, PAN-OS 8.0 only) Auto.

`server` supports the following arguments:

* `name` - (Required) The name.
* `server` - (Required) Server hostname or IP address.
* `secret` - (Required) Shared secret for TACACS+ communication.
* `port` - (int) TACACS+ server port number (default: `49`).
