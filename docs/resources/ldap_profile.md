---
page_title: "panos: panos_ldap_profile"
subcategory: "Device"
---

# panos_ldap_profile

Manages a LDAP profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_ldap_profiles" "example" {
    name = "fromTerraform"
    base_dn = "baseDn"
    bind_dn = "bindDn"
    password = "secret"
    bind_timeout = 5
    search_timeout = 7
    retry_interval = 120
    server {
        name = "first"
        server = "first.example.com"
    }
    server {
        name = "second"
        server = "192.168.0.5"
        port = 23430
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
* `ldap_type` - LDAP type.  Valid values are `"active-directory"`, `"e-directory"`, `"sun"`, or `"other"` (default).
* `ssl` - (bool) SSL (default: `true`).
* `verify_server_certificate` - (bool) Verify server certificate for SSL sessions.
* `disabled` - (bool) Disable this profile.
* `base_dn` - Default base distinguished name (DN) to use for searches.
* `bind_dn` - Bind distinguished name.
* `password` - (Required) Bind password.
* `search_timeout` - (int) Number of seconds to wait for performing searches (default: `30`).
* `bind_timeout` - (int) Number of seconds to use for connecting to servers (default: `30`).
* `retry_interval` - (int) Interval (in seconds) for reconnecting LDAP server (default: `60`).
* `server` - List of server specs, as defined below.

`server` supports the following arguments:

* `name` - (Required) The name.
* `server` - (Required) Server hostname or IP address.
* `port` - (int) LDAP server port number (default: `389`).
