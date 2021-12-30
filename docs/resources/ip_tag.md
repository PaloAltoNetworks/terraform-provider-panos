---
page_title: "panos: panos_ip_tag"
subcategory: "User-ID"
---

# panos_ip_tag

Manages a specific set of tags on a single IP address.

This resource only manages the given tags for the given IP address.  Any
other tags associated with the IP address are left as-is.


## PAN-OS

NGFW and Panorama (8.0+)


## Example Usage

```hcl
resource "panos_ip_tag" "example1" {
    ip = "10.2.3.4"
    tags = [
        "tag1",
        "tag2",
    ]
}

# It is safe to have multiple resources target the same IP.
resource "panos_ip_tag" "example2" {
    ip = "10.2.3.4"
    tags = [
        "tag3",
    ]
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - The vsys location (default: `vsys1`).
* `ip` - (Required) The IP address.
* `tags` - (list) List of tags.
