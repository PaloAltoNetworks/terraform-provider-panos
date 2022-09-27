---
page_title: "panos: panos_license_api_key"
subcategory: "Device"
---

# panos_license_api_key

This resource manages the licensing API key, which is necessary to delicense
the PAN-OS firewall.

This resource's `retain_key` param is a Terraform side configuration only.  In
order for the firewall to delicense itself, the licensing API key must be
present.  This means that either the `panos_licensing` resource must use
`depends_on` and depend on this resource, or you must set the `retain_key`
param to `true`.  As there is no harm in leaving the licensing API key on the
PAN-OS firewall, it is recommended that `retain_key` be set to `true`.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_license_api_key" "example" {
    key = "secret"
    retain_key = true

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) The licensing API key.
* `retain_key` - (Optional) Set to `true` to retain the licensing API key
  even after the deletion of this resource (recommended).
