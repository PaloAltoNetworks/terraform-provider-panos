---
layout: "panos"
page_title: "panos: panos_panorama_application_object"
description: |-
  Manages Panorama application objects.
---

# panos_panorama_application_object

This resource allows you to add/update/delete Panorama application objects.


## Import Name

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_application_object" "example" {
    name = "myApp"
    description = "made by terraform"
    category = "media"
    subcategory = "gaming"
    technology = "browser-based"
    defaults {
        port {
            ports = [
                "udp/dynamic",
            ]
        }
    }
    risk = 4
    scanning {
        viruses = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The object's name.
* `device_group` - (Optional) The device group (default: `shared`)
* `defaults` - (Optional) The application's defaults spec (defined below).  To have
  a "defaults" of `None`, omit this section.
* `category` - (Required) The category.
* `subcategory` - (Required) The subcategory.
* `technology` - (Required) The technology.
* `description` - (Optional) The object's description.
* `timeout_settings` - (Optional) The timeout spec (defined below).
* `risk` - (Optional, int) The risk (default: 1).
* `parent_app` - (Optional) The parent application.
* `able_to_file_transfer` - (Optional, bool) Able to file transfer.
* `excessive_bandwidth` - (Optional, bool) Excessive bandwidth use.
* `tunnels_other_applications` - (Optional, bool) This application tunnels other apps.
* `has_known_vulnerability` - (Optional, bool) Has known vulnerabilities.
* `used_by_malware` - (Optional, bool) App is used by malware.
* `evasive_behavior` - (Optional, bool) App is evasive.
* `pervasive_use` - (Optional, bool) App is pervasive.
* `prone_to_misuse` - (Optional, bool) Prone to misuse.
* `continue_scanning_for_other_applications` - (Optional, bool) Continue scanning for
  other applications.
* `scanning` - The scanning spec (defined below).
* `alg_disable_capability` - (Optional) The alg disable capability.
* `no_app_id_caching` - (Optional, bool) No appid caching.

`defaults` supports the following arguments:

* `port` - (Optional) The port spec (defined below)
* `ip_protocol` - (Optional) The ip protocol spec (defined below)
* `icmp` - (Optional) The ICMP spec (defined below)
* `icmp6` - (Optional) The ICMP6 spec (defined below)

`defaults.port` supports the following arguments:

* `ports` - (Required) List of ports.

`defaults.ip_protocol` supports the following arguments:

* `value` - (Required, int) The IP protocol value.

`defaults.icmp` supports the following arguments:

* `type` - (Required, int) The type.
* `code` - (Optional, int) The code.

`defaults.icmp6` supports the following arguments:

* `type` - (Required, int) The type.
* `code` - (Optional, int) The code.

`timeout_settings` supports the following arguments:

* `timeout` - (Optional, int) The timeout.
* `tcp_timeout` - (Optional, int) TCP timeout.
* `udp_timeout` - (Optional, int) UDP timeout.
* `tcp_half_closed` - (Optional, int) TCP half closed timeout.
* `tcp_time_wait` - (Optional, int) TCP time wait timeout.

`scanning` supports the following arguments:

* `file_types` - (Optional, bool) File type scanning.
* `viruses` - (Optional, bool) Virus scanning.
* `data_patterns` - (Optional, bool) Data pattern scanning.
