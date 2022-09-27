---
page_title: "panos: panos_url_filtering_security_profile"
subcategory: "Objects"
---

# panos_url_filtering_security_profile

Manages URL filtering security profiles.


## Import Name

NGFW:

```shell
<vsys>:<name>
```

Panorama:

```shell
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_url_filtering_security_profile" "example" {
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

    lifecycle {
        create_before_destroy = true
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
* `description` - The description.
* `dynamic_url` - (bool, removed in PAN-OS 9.0) Dynamic URL.
* `expired_license_action` - (bool, removed in PAN-OS 9.0) Expired license action.
* `block_list_action` - (removed in PAN-OS 9.0) Block list action.
* `block_list` - (list, removed in PAN-OS 9.0) Block list.
* `allow_list` - (list, removed in PAN-OS 9.0) Allow list.
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
* `ucd_mode` - (PAN-OS 8.0+) User credential detection mode.  Valid values are
  `disabled` (default), `ip-user`, `domain-credentials`, or `group-mapping`.
* `ucd_mode_group_mapping` - (`ucd_mode`=`group-mapping`, PAN-OS 8.0+) User
  credential detection: the group mapping settings.
* `ucd_log_severity` - (Optional, but Required for PAN-OS 8.0+) User credential
  detection: valid username detected log severity.
* `ucd_allow_categories` - (list, PAN-OS 8.0+) Categories allowed with user
  credential submission.
* `ucd_alert_categories` - (list, PAN-OS 8.0+) Categories alerted on with
  user credential submission.
* `ucd_block_categories` - (list, PAN-OS 8.0+) Categories blocked with
  user credential submission.
* `ucd_continue_categories` - (list, PAN-OS 8.0+) Categories continued with
  user credential submission.
* `http_header_insertion` - (repeatable, PAN-OS 8.1+) List of HTTP header
  insertion specs, as defined below.
* `machine_learning_model` - (repeatable, PAN-OS 10.0+) List of machine learning
  specs, as defined below.
* `machine_learning_exceptions` - (list) List of machine learning exceptions.

`http_header_insertion` supports the following arguments:

* `name` - (Required) Name.
* `type` - Header type.
* `domains` - (list) Header domains.
* `http_header` - (repeatable) List of HTTP header specs, as defined below.

`http_header_insertion.http_header` supports the following arguments:

* `name` - (Computed attribute) Name.
* `header` - (Required) The header.
* `value` - (Required) The value of the header.
* `log` - (bool) Logging of this header.

`machine_learning_model` supports the following arguments:

* `model` - (Required) Machine learning model.
* `action` - (Required) Model action.
