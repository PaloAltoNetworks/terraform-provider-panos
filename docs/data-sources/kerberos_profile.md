---
page_title: "panos: panos_kerberos_profile"
subcategory: "Device"
---

# panos_kerberos_profile

Gets information on a Kerberos profile.


## Minimum PAN-OS Version

7.0


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_kerberos_profiles" "example" {
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

* `admin_use_only` - (bool) Administrator use only.
* `server` - List of server specs, as defined below.

`server` supports the following arguments:

* `name` - The name.
* `server` - Server hostname or IP address.
* `port` - (int) Kerberos server port number.
