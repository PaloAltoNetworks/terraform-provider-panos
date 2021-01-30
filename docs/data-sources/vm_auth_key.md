---
page_title: "panos: panos_vm_auth_key"
subcategory: "Panorama"
---


# panos_vm_auth_key

Gets info on VM auth keys on Panorama.

**NOTE:** This is for Panorama only.


## Example Usage

```hcl
data "panos_vm_auth_key" "example" {}
```


## Attribute Reference

The following attributes are supported.

* `total` - (int) Total number of entries.
* `entries` - List of entry structs, as defined below.

`entries` supports the following attributes:

* `auth_key` - The VM auth key.
* `expiry` - The expiry time as reported by PAN-OS
* `valid` - (bool) If the VM auth key is still valid or not.
