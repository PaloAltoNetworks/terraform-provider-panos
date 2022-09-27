---
page_title: "panos: panos_http_server_profile"
subcategory: "Device"
---

# panos_http_server_profile

This resource allows you to add/update/delete HTTP server profiles.


## Minimum PAN-OS Version

7.1


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_http_server_profile" "example" {
    name = "myProfile"
    url_format {
        name = "my url format"
        uri_format = "/api/incident/url"
        headers = {
            "Content-Type": "text/plain",
        }
    }
    http_server {
        name = "myServer"
        address = "siem.example.com"
        certificate_profile = data.panos_system_info.x.version_major >= 9 ? "None" : ""
        tls_version = data.panos_system_info.x.version_major >= 9 ? "1.2" : ""
    }

    lifecycle {
        create_before_destroy = true
    }
}

data "panos_system_info" "x" {}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys (default: `shared`).
* `name` - (Required) The group's name.
* `tag_registration` - (Optional, bool) Perform tag registration.
* `config_format` - (Optional) A format folder spec for config (defined below).
* `system_format` - (Optional) A format folder spec for system (defined below).
* `threat_format` - (Optional) A format folder spec for threat (defined below).
* `traffic_format` - (Optional) A format folder spec for traffic (defined below).
* `hip_match_format` - (Optional) A format folder spec for HIP match (defined below).
* `url_format` - (Optional) A format folder spec for url (defined below).
* `data_format` - (Optional) A format folder spec for data (defined below).
* `wildfire_format` - (Optional) A format folder spec for wildfire (defined below).
* `tunnel_format` - (Optional) A format folder spec for tunnel (defined below).
* `user_id_format` - (Optional) A format folder spec for user ID (defined below).
* `gtp_format` - (Optional) A format folder spec for gtp (defined below).
* `auth_format` - (Optional) A format folder spec for auth (defined below).
* `sctp_format` - (Optional, PAN-OS 8.1+) A format folder spec for sctp (defined below).
* `iptag_format` - (Optional, PAN-OS 9.0+) A format folder spec for iptag (defined below).
* `http_server` - (Optional, repeatable) A server spec (defined below).

All format folders support the following arguments:

* `name` - (Optional) The name.
* `uri_format` - (Optional) The URI format.
* `payload` - (Optional) The payload.
* `headers` - (Optional, map) A map of HTTP headers and their values.
* `params` - (Optional, map) A map of HTTP params and their values.

`http_server` supports the following arguments:

* `name` - (Required) The server name.
* `address` - (Required) The server address.
* `protocol` - (Optional) The protocol.  Valid values are `HTTPS` (default)
  or `HTTP`.
* `port` - (Optional, int) The port number (default: 443).
* `http_method` - (Optional) The HTTP method (default: `POST`).
* `username` - (Optional) The username.
* `password` - (Optional) The password.
* `tls_version` - (Optional) The TLS version.
* `certificate_profile` - (Optional) The certificate profile.
