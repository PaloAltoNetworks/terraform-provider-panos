---
page_title: "panos: panos_monitor_profile"
subcategory: "Network"
---

# panos_monitor_profile

This resource allows you to add/update/delete monitor profiles.


## Minimum PAN-OS Version

7.1


## PAN-OS

NGFW


## Import Name

```shell
<name>
```

## Example Usage

```hcl
resource "panos_monitor_profile" "example" {
    name = "myProfile"
    interval = 5
    threshold = 3

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The monitor profile name.
* `interval` - (Optional, int) The probing interval in seconds.
* `threshold` - (Optional, int) The number of failed probes to determine that
  the tunnel is down.
* `action` - (Optional) Action triggered when tunnel's status changes.  Valid values
  are `wait-recover` (default) or `fail-over`.
