---
page_title: "panos: panos_snmptrap_server_profile"
subcategory: "Device"
---

# panos_snmptrap_server_profile

This resource allows you to add/update/delete snmptrap server profiles.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_snmptrap_server_profile" "example" {
    name = "myProfile"
    v2c_server {
        name = "first server"
        manager = "snmp1.example.com"
        community = "public"
    }

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys (default: `shared`).
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
