---
page_title: "panos: panos_data_filtering_security_profile"
subcategory: "Objects"
---

# panos_data_filtering_security_profile

Gets information on a data filtering security profile.


## Example Usage

```hcl
data "panos_data_filtering_security_profile" "example" {
    name = panos_data_filtering_security_profile.x.name
}

resource "panos_data_filtering_security_profile" "x" {
    name = "example"
    description = "made by Terraform"
    rule {
        data_pattern = panos_custom_data_pattern_object.my_custom_obj.name
        applications = ["any"]
        file_types = ["any"]
    }

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_custom_data_pattern_object" "my_custom_obj" {
    name = "myobj"
    description = "for my data filtering security profile"
    type = "regex"
    regex {
        name = "my regex"
        file_types = ["any"]
        regex = "this is my regex"
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


## Attribute Reference

The following attributes are supported:

* `description` - The description.
* `data_capture` - (bool) Data capture.
* `rule` - (repeatable) Rule list spec, as defined below.

`rule` supports the following arguments:

* `name` - Name.
* `data_pattern` - The data pattern name.
* `applications` - (list) List of applications.
* `file_types` - (list) List of file types.
* `direction` - Direction.
* `alert_threshold` - (int) Alert threshold.
* `block_threshold` - (int) Block threshold.
* `log_severity` - Log severity.
