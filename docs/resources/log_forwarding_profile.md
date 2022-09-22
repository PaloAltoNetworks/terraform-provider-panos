---
page_title: "panos: panos_log_forwarding_profile"
subcategory: "Objects"
---

# panos_log_forwarding_profile

This resource allows you to add/update/delete log forwarding profiles.


## Minimum PAN-OS Version

8.0


## PAN-OS

NGFW


## Import Name

```shell
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_log_forwarding_profile" "example" {
    name = "myProfile"
    description = "made by Terraform"
    match_list {
        name = "myMatchList"
        log_type = "url"
        http_server_profiles = [
            panos_http_server_profile.h1.name,
            panos_http_server_profile.h2.name,
        ]
        action {
            name = "tagging int"
            tagging_integration {
                timeout = 5
                local_registration {
                    tags = [
                        panos_administrative_tag.t.name,
                    ]
                }
            }
        }
        action {
            name = "azure int"
            azure_integration { }
        }
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_http_server_profile" "h1" {
    name = "h1"
    http_server {
        name = "h1"
        address = "h1.example.com"
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_http_server_profile" "h2" {
    name = "h2"
    http_server {
        name = "h2"
        address = "h2.example.com"
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_administrative_tag" "t" {
    name = "myTag"
    color = "color12"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The vsys (default: `vsys1`).
* `name` - (Required) The group's name.
* `description` - (Optional) The description.
* `enhanced_logging` - (Optional, bool, PAN-OS 8.1+) Set to `true` to enable enhanced logging.
* `match_list` - (Optional, repeatable) The match list spec (defined below).

`match_list` supports the following arguments:

* `name` - (Required) The name.
* `description` - (Optional) The description.
* `log_type` - (Optional) The log type.  Valid values are `traffic` (default),
  `threat`, `wildfire`, `url`, `data`, `gtp`, `tunnel`, `auth`, `sctp`, or `decryption`.
* `filter` - (Optional) The filter (default: `All Logs`).
* `send_to_panorama` - (Optional, bool) Set to `true` to send to Panorama.
* `snmptrap_server_profiles` - (Optional) List of SNMP server profiles.
* `email_server_profiles` - (Optional) List of email server profiles.
* `syslog_server_profiles` - (Optional) List of syslog server profiles.
* `http_server_profiles` - (Optional) List of http server profiles.
* `action` - (Optional, repeatable) Match list action spec (defined below).

`match_list.action` supports the following arguments:

* `name` - (Required) The name.
* `azure_integration` - (Optional) The Azure integration spec (defined
  below).  Mutually exclusive with `tagging_integration`.
* `tagging_integration` - (Optional) The tagging integration spec (defined
  below).  Mutually exclusive with `azure_integration`.

`match_list.action.azure_integration` supports the following arguments:

* `azure_integration` - (Optional, bool) This param defaults to `true` and should
  not be changed.

`match_list.action.tagging_integration` supports the following arguments:

* `action` - (Optional) The action.  Valid values are `add-tag` (default) or
  `remove-tag`.
* `target` - (Optional) The target.  Valid values are `source-address` (default)
  or `destination-address`.
* `timeout` - (Optional, int) The timeout.
* `local_registration` - (Optional) The local registration spec (defined below).
  Only one of `local_registration`, `remote_registration`, or `panorama_registration`
  should be defined.
* `remote_registration` - (Optional) The remote registration spec (defined below).
  Only one of `local_registration`, `remote_registration`, or `panorama_registration`
  should be defined.
* `panorama_registration` - (Optional) The panorama registration spec (defined below).
  Only one of `local_registration`, `remote_registration`, or `panorama_registration`
  should be defined.

`match_list.action.tagging_integration.local_registration` supports the
following arguments:

* `tags` - (Required) List of administrative tags.

`match_list.action.tagging_integration.remote_registration` supports the
following arguments:

* `tags` - (Required) List of administrative tags.
* `http_profile` - (Required) The HTTP profile.

`match_list.action.tagging_integration.panorama_registration` supports the
following arguments:

* `tags` - (Required) List of administrative tags.
