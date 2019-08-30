---
layout: "panos"
page_title: "panos: panos_syslog_server_profile"
description: |-
  Manages syslog server profiles.
---

# panos_syslog_server_profile

This resource allows you to add/update/delete syslog server profiles.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_syslog_server_profile" "example" {
    name = "myProfile"
    threat_format = "$serial $severity"
    syslog_server {
        name = "my-server"
        server = "syslog.example.com"
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
