---
layout: "panos"
page_title: "panos: panos_licensing"
sidebar_current: "docs-panos-resource-licensing"
description: |-
  Manages PAN-OS licensing.
---

# panos_licensing

This resource manages the licenses installed on the PAN-OS firewall.

## Example Usage

```hcl
resource "panos_licensing" "example" {
    auth_codes = ["code1", "code2"]
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
