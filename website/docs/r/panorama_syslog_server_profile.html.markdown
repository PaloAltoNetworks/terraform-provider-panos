---
layout: "panos"
page_title: "panos: panos_panorama_syslog_server_profile"
description: |-
  Manages Panorama syslog server profiles.
---

# panos_panorama_syslog_server_profile

This resource allows you to add/update/delete Panorama syslog server profiles.


## Import Name

```
<template>:<template_stack>:<vsys>:<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_syslog_server_profile" "example" {
    device_group = "shared"
    name = "myProfile"
    threat_format = "$serial $severity"
    syslog_server {
        name = "my-server"
        server = "syslog.example.com"
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
* `syslog_server` - (Required, repeatable) The server spec (defined below).

`syslog_server` supports the following arguments:

* `name` - (Required) Server name.
* `server` - (Required) The hostname.
* `transport` - (Optional) The transport.  Valid values are `UDP` (default),
  `TCP`, or `SSL`.
* `port` - (Optional, int) The port number (default: 514).
* `syslog_format` - (Optional) The syslog format.  Valid values are `BSD`
  (default) or `IETF`.
* `facility` - (Optional) The syslog facility.  Valid values are `LOG_USER`
  (default), `LOG_LOCAL0`, `LOG_LOCAL1`, `LOG_LOCAL2`, `LOG_LOCAL3`,
  `LOG_LOCAL4`, `LOG_LOCAL5`, `LOG_LOCAL6`, or `LOG_LOCAL7`.
