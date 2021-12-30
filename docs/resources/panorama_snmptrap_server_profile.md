---
page_title: "panos: panos_panorama_snmptrap_server_profile"
subcategory: "Device"
---

# panos_panorama_snmptrap_server_profile

This resource allows you to add/update/delete Panorama snmptrap server profiles.


## PAN-OS

Panorama


## Example Usage

```hcl
resource "panos_panorama_snmptrap_server_profile" "example" {
    device_group = "shared"
    name = "myProfile"
    v2c_server {
        name = "first server"
        manager = "snmp1.example.com"
        community = "public"
    }
}
```

## Argument Reference

When creating this profile, there are a few options:

* on the Panorama (in the `shared` device group)
* in a vsys in a template
* in a vsys in a template stack

The following arguments are supported:

* `template` - (Optional) The template location.  Mutually exclusive with
  `template_stack` and `device_group`.
* `template_stack` - (Optional) The template stack location.  Mutually exclusive
  with `template` and `device_group`.
* `device_group` - (Optional) The device group location.  Mutually exclusive with
  `template` and `template_stack`.
* `vsys` - (Optional) The vsys.  This will likely be `shared`, and it should be
  defined if you specified either `template` or `template_stack`.
* `name` - (Required) The group's name.
* `v2c_server` - (Optional, repeatable) A v2c server (defined below).
* `v3_server` - (Optional, repeatable) A v3 server (defined below).

`v2c_server` supports the following arguments:

* `name` - (Required) The server name.
* `manager` - (Required) The hostname.
* `community` - (Required) The SNMP community.

`v3_server` supports the following arguments:

* `name` - (Required) The server name.
* `manager` - (Required) The hostname.
* `user` - (Required) Username.
* `engine_id` - (Optional) The engine ID.
* `auth_password` - (Required) SNMP auth password.
* `priv_password` - (Required) SNMP priv password.
