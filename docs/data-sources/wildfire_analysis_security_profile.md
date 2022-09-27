---
page_title: "panos: panos_wildfire_analysis_security_profile"
subcategory: "Objects"
---

-> **NOTE:**  Minimum PAN-OS version: 7.0


# panos_wildfire_analysis_security_profile

Gets information on wildfire analysis security profiles.


## Import Name

NGFW:

```
<vsys>:<name>
```

Panorama:

```
<device_group>:<name>
```


## Example Usage

```hcl
data "panos_wildfire_analysis_security_profile" "example" {
    name = panos_wildfire_analysis_security_profile.x.name
}

resource "panos_wildfire_analysis_security_profile" "x" {
    name = "example"
    description = "made by Terraform"
    rule {
        name = "foo"
        applications = ["pop3"]
        file_types = ["pdf"]
    }

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `rule` - (repeatable) Rule list spec, as defined below.

`rule` supports the following arguments:

* `name` - (Required) Name.
* `applications` - (list) List of applications.
* `file_types` - (list) List of file types.
* `direction` - Direction.  Valid values are `both` (default),
  `upload`, or `download`.
* `analysis` - Analysis setting.  Valid values are `public-cloud` (default)
  or `private-cloud`.
