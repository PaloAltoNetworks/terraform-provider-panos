---
page_title: "panos: panos_dos_protection_profile"
subcategory: "Objects"
---

# panos_dos_protection_profile

Gets DOS protection security profile info.


## Example Usage

```hcl
data "panos_dos_protection_profile" "example" {
    name = panos_dos_protection_profile.x.name
}

resource "panos_dos_protection_profile" "x"
    name = "example"
    description = "made by Terraform"
    syn {
        enable = True
        action = "red"
        alarm_rate = 777
        activate_rate = 888
        max_rate = 999
        block_duration = 42
    }

    lifecycle {
        create_before_destroy = true
    }
}
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
* `type` - The profile type.
* `enable_sessions_protections` - (bool) Enable sessions protections.
* `max_concurrent_sessions` - (int) Max concurrent sessions.
* `syn` - Optional syn spec, as defined below.
* `udp` - Optional protection spec, as defined below.
* `icmp` - Optional ICMP spec, as defined below.
* `icmpv6` - Optional ICMPv6 spec, as defined below.
* `other` - Optional other IP flood protection spec, as defined below.

`syn` supports the following attributes:

* `enable` - (bool) Enable
* `action` - SYN protection action.
* `alarm_rate` - (int) Alarm rate.
* `activate_rate` - (int) Activate rate.
* `max_rate` - (int) Max rate.
* `block_duration` - (int) Block duration.

`udp`, `icmp`, `icmpv6`, and `other` all support the following attributes:

* `enable` - (bool) Enable
* `alarm_rate` - (int) Alarm rate.
* `activate_rate` - (int) Activate rate.
* `max_rate` - (int) Max rate.
* `block_duration` - (int) Block duration.
