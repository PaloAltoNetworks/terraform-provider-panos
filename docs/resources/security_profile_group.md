---
page_title: "panos: panos_security_profile_group"
subcategory: "Objects"
---

# panos_security_profile_group

This resource allows you to add/update/delete a security profile group.


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<device_group>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_security_profile_group" "example" {
    name = "myGroup"
    antivirus_profile = "default"
    anti_spyware_profile = "anti-spyware1"

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `device_group` - The device group (default: `shared`)


NGFW:

* `vsys` - The vsys (default: `vsys1`).


The following arguments are supported:

* `name` - (Required) The name.
* `antivirus_profile` - The AV profile name.
* `anti_spyware_profile` - Anti Spyware profile name.
* `vulnerability_profile` - Vulnerability profile name.
* `url_filtering_profile` - URL filtering profile name.
* `file_blocking_profile` - File blocking profile name.
* `data_filtering_profile` - Data filtering profile name.
* `wildfire_analysis_profile` - Wildfire analysis profile name.
* `gtp_profile` - GTP profile name.
* `sctp_profile` - SCTP profile name.
