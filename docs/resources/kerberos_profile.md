---
page_title: "panos: panos_kerberos_profile"
subcategory: "Device"
---

# panos_kerberos_profile

Manages a Kerberos profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_kerberos_profiles" "example" {
    name = "fromTerraform"
    admin_use_only = true
    server {
        name = "server1"
        server = "kerberos1.example.com"
    }
    server {
        name = "server2"
        server = "music.example.com"
        port = 12345
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
* `server` - List of server specs, as defined below.

`server` supports the following arguments:

* `name` - (Required) The name.
* `server` - (Required) Server hostname or IP address.
* `port` - (int) Kerberos server port number (default: `88`).
