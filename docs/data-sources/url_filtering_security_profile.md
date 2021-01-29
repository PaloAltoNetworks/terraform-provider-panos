---
page_title: "panos: panos_url_filtering_security_profile"
subcategory: "Objects"
---

# panos_url_filtering_security_profile

Gets information on URL filtering security profiles.


## Example Usage

```hcl
data "panos_url_filtering_security_profile" "example" {
    name = panos_url_filtering_security_profile.x.name
}

resource "panos_url_filtering_security_profile" "x" {
    name = "example"
    description = "made by terraform"
    ucd_mode = "disabled"
    ucd_log_severity = "${
        data.panos_system_info.x.version_major > 8 ? "medium" : ""
    }"
    log_container_page_only = true
    log_http_header_xff = true
    log_http_header_referer = true
    log_http_header_user_agent = true
    http_header_insertion {
        name = "doublelift"
        type = "Custom"
        domains = [
            "b.example.com",
            "a.example.com",
            "c.example.com",
        ]
        http_header {
            header = "X-First-Header"
            value = "alpha"
        }
        http_header {
            header = "X-Second-Header"
            value = "beta"
        }
    }
}

data "panos_system_info" "x" {}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `description` - The description.
* `dynamic_url` - (bool) Dynamic URL.
* `expired_license_action` - (bool) Expired license action.
* `block_list_action` - Block list action.
* `block_list` - (list) Block list.
* `allow_list` - (list) Allow list.
* `allow_categories` - (list) List of categories to allow.
* `alert_categories` - (list) List of categories to alert.
* `block_categories` - (list) List of categories to block.
* `continue_categories` - (list) List of categories to continue.
* `override_categories` - (list) List of categories to override.
* `track_container_page` - (bool) Track the container page.
* `log_container_page_only` - (bool) Log container page only.
* `safe_search_enforcement` - (bool) Safe search enforcement.
* `log_http_header_xff` - (bool) HTTP header logging: X-Forwarded-For.
* `log_http_header_user_agent` - (bool) HTTP header logging: User-Agent.
* `log_http_header_referer` - (bool) HTTP header logging: Referer.
* `ucd_mode` - User credential detection mode.
* `ucd_mode_group_mapping` - (`ucd_mode`=`group-mapping`) User
  credential detection: the group mapping settings.
* `ucd_log_severity` - User credential detection: valid username detected log severity.
* `ucd_allow_categories` - (list) Categories allowed with user
  credential submission.
* `ucd_alert_categories` - (list) Categories alerted on with
  user credential submission.
* `ucd_block_categories` - (list) Categories blocked with
  user credential submission.
* `ucd_continue_categories` - (list) Categories continued with
  user credential submission.
* `http_header_insertion` - (repeatable) List of HTTP header
  insertion specs, as defined below.
* `machine_learning_model` - (repeatable) List of machine learning
  specs, as defined below.
* `machine_learning_exceptions` - (list) List of machine learning exceptions.

`http_header_insertion` supports the following attributes:

* `name` - Name.
* `type` - Header type.
* `domains` - (list) Header domains.
* `http_header` - (repeatable) List of HTTP header specs, as defined below.

`http_header_insertion.http_header` supports the following attributes:

* `name` - Name.
* `header` - The header.
* `value` - The value of the header.
* `log` - (bool) Logging of this header.

`machine_learning_model` supports the following arguments:

* `model` - Machine learning model.
* `action` - Model action.
