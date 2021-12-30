---
page_title: "panos: panos_ip_tag"
subcategory: "User-ID"
---

# panos_ip_tag

Gets information on IP tags.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
# Get all IPs and all tags.
data "panos_ip_tag" "example1" {}

# Get all tags for a single IP.
data "panos_ip_tag" "example2" {
    ip = "10.2.3.4"
}

# Get all IPs with tag "foo".
data "panos_ip_tag" "example3" {
    tag = "foo"
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - The vsys location (default: `vsys1`).
* `ip` - Filter on a specific IP address.
* `tag` - Filter on a specific tag.


## Attribute Reference

The following attributes are supported:

* `total` - (int) Total number of entries.
* `entries` - List of entries specs, as defined below.

`entries` supports the following attributes:

* `ip` - The IP address.
* `tags` - (list) List of tags.
