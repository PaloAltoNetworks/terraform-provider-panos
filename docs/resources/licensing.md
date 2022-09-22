---
page_title: "panos: panos_licensing"
subcategory: "Device"
---

# panos_licensing

This resource manages the licenses installed on the PAN-OS firewall.

Installing the standard auth code for the standard PAN-OS license key for the
firewall causes the firewall to reboot.  Thus it is recommended that you use
this resource in a separate step of your overall firewall provisioning, as
using this resource will cause the firewall to be temporarily inaccessible.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_licensing" "example" {
    auth_codes = [
        "code1",
        "code2",
    ]

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `auth_codes` - (Required) The list of auth codes to install.
* `delicense` - (Optional, bool) Leave as `true` if you want to delicense
  the firewall when this resource is removed, otherwise set to `false` to
  prevent firewall delicensing.  Delicensing requires that the licensing
  API key has been installed.
* `mode` - (Optional) For `delicense` of `true`, the type of delicensing to
  perform.  Right now, only `auto` is supported (no manual delicensing).

## Attribute Reference

The following attributes are available after read operations:

* `licenses` - List of licenses.

Licenses have the following attributes:

* `feature` - The feature name.
* `description` - License description.
* `serial` - The serial number.
* `issued` - When the license was issued.
* `expires` - When the license expires.
* `expired` - If the license has expired or not.
* `auth_code` - Associated auth code (if applicable).
