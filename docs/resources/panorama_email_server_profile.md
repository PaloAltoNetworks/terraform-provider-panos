---
page_title: "panos: panos_panorama_email_server_profile"
subcategory: "Device"
---

# panos_panorama_email_server_profile

This resource allows you to add/update/delete Panorama email server profiles.


## PAN-OS

Panorama


## Import Name

```
<template>:<template_stack>:<vsys>:<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_email_server_profile" "example" {
    name = "myProfile"
    device_group = "shared"
    threat_format = "$serial $severity"
    email_server {
        name = "my-server"
        display_name = "foobar"
        from_email = "source@example.com"
        to_email = "alerts@example.com"
        email_gateway = "mail.example.com"
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
* `config_format` - (Optional) Config format.
* `system_format` - (Optional) System format.
* `threat_format` - (Optional) Threat format.
* `traffic_format` - (Optional) Traffic format.
* `hip_match_format` - (Optional) HIP match format.
* `url_format` - (Optional) URL format.
* `data_format` - (Optional) Data format.
* `wildfire_format` - (Optional) Wildfire format.
* `tunnel_format` - (Optional) Tunnel format.
* `user_id_format` - (Optional) UserID format.
* `gtp_format` - (Optional) GTP format.
* `auth_format` - (Optional) Auth format.
* `sctp_format` - (Optional) SCTP format.
* `iptag_format` - (Optional) IP tag format.
* `escaped_characters` - (Optional) The escaped characters (as a string).
* `escape_character` - (Optional) The escape character.
* `email_server` - (Required, repeatable) The server spec (defined below).

`email_server` supports the following arguments:

* `name` - (Required) Server name.
* `display_name` - (Optional) The display name.
* `from_email` - (Required) From email address.
* `to_email` - (Required) To email address.
* `also_to_email` - (Optional) Secondary to email address.
* `email_gateway` - (Required) The email server.
* `email_server` - (Required, repeatable) The server spec (defined below).
