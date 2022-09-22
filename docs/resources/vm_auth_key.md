---
page_title: "panos: panos_vm_auth_key"
subcategory: "Panorama"
---


# panos_vm_auth_key

Creates a VM auth key you can use to bootstrap a VM NGFW.

**NOTE:** This is for Panorama only.


## Example Usage

```hcl
# Basic usage.
resource "panos_vm_auth_key" "example1" {
    hours = 24

    lifecycle {
        create_before_destroy = true
    }
}
```

```hcl
# Using time_rotating to keep a VM auth key always valid.
resource "panos_vm_auth_key" "example2" {
    hours = time_rotating.tr.rotation_hours

    keepers = {
        rotate = time_rotating.tr.rotation_rfc3339
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "time_rotating" "tr" {
    rotation_hours = 24
}
```


## Argument Reference

The following arguments are supported:

* `hours` - (int) The VM auth key lifetime in hours.
* `keepers` - (map) Arbitrary map of values that, when changed, will trigger a new key to be generated.


## Attribute Reference

The following attributes are supported.

* `auth_key` - The bootstrap VM auth key.
* `expiry` - The date as returned from PAN-OS for when the auth key expires.
* `expiration` - The expiration time as a RFC 3339 string.
* `valid` - (bool) If the auth key is still valid based on the lifetime given.
