---
page_title: "panos: panos_security_profile_group"
subcategory: "Objects"
---

# panos_security_profile_group

Retrieve information on the specified security profile group.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_security_profile_group" "example" {
    name = "myGroup"
}
```


## Argument Reference

Panorama:

* `device_group` - (Optional) The device group (default: `shared`)


NGFW:

* `vsys` - (Optional) The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `antivirus_profile` - The AV profile name.
* `anti_spyware_profile` - Anti Spyware profile name.
* `vulnerability_profile` - Vulnerability profile name.
* `url_filtering_profile` - URL filtering profile name.
* `file_blocking_profile` - File blocking profile name.
* `data_filtering_profile` - Data filtering profile name.
* `wildfire_analysis_profile` - Wildfire analysis profile name.
* `gtp_profile` - GTP profile name.
* `sctp_profile` - SCTP profile name.
