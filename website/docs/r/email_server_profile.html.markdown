---
layout: "panos"
page_title: "panos: panos_email_server_profile"
description: |-
  Manages email server profiles.
---

# panos_email_server_profile

This resource allows you to add/update/delete email server profiles.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_email_server_profile" "example" {
    name = "myProfile"
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

The following arguments are supported:

* `vsys` - (Optional) The vsys (default: `shared`).
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
