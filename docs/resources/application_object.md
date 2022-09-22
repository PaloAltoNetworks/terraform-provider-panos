---
page_title: "panos: panos_application_object"
subcategory: "Objects"
---

# panos_application_object

This resource allows you to add/update/delete application objects.


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<device_group>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_application_object" "example" {
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

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

Panorama:

* `device_group` - The device group (default: `shared`).

NGFW:

* `vsys` - The object's vsys (default: `vsys1`).

The following arguments are supported:

* `name` - (Required) The object's name.
* `defaults` - The application's defaults spec (defined below).  To have
  a "defaults" of `None`, omit this section.
* `category` - (Required) The category.
* `subcategory` - (Required) The subcategory.
* `technology` - (Required) The technology.
* `description` - The object's description.
* `timeout_settings` - The timeout spec (defined below).
* `risk` - (int) The risk (default: 1).
* `parent_app` - The parent application.
* `able_to_file_transfer` - (bool) Able to file transfer.
* `excessive_bandwidth` - (bool) Excessive bandwidth use.
* `tunnels_other_applications` - (bool) This application tunnels other apps.
* `has_known_vulnerability` - (bool) Has known vulnerabilities.
* `used_by_malware` - (bool) App is used by malware.
* `evasive_behavior` - (bool) App is evasive.
* `pervasive_use` - (bool) App is pervasive.
* `prone_to_misuse` - (bool) Prone to misuse.
* `continue_scanning_for_other_applications` - (bool) Continue scanning for
  other applications.
* `scanning` - The scanning spec (defined below).
* `alg_disable_capability` - The alg disable capability.
* `no_app_id_caching` - (bool) No appid caching.

`defaults` supports the following arguments:

* `port` - The port spec (defined below)
* `ip_protocol` - The ip protocol spec (defined below)
* `icmp` - The ICMP spec (defined below)
* `icmp6` - The ICMP6 spec (defined below)

`defaults.port` supports the following arguments:

* `ports` - (Required) List of ports.

`defaults.ip_protocol` supports the following arguments:

* `value` - (Required) The IP protocol value.

`defaults.icmp` supports the following arguments:

* `type` - (Required, int) The type.
* `code` - (int) The code.

`defaults.icmp6` supports the following arguments:

* `type` - (Required, int) The type.
* `code` - (int) The code.

`timeout_settings` supports the following arguments:

* `timeout` - (int) The timeout.
* `tcp_timeout` - (int) TCP timeout.
* `udp_timeout` - (int) UDP timeout.
* `tcp_half_closed` - (int) TCP half closed timeout.
* `tcp_time_wait` - (int) TCP time wait timeout.

`scanning` supports the following arguments:

* `file_types` - (bool) File type scanning.
* `viruses` - (bool) Virus scanning.
* `data_patterns` - (bool) Data pattern scanning.
