---
page_title: "panos: panos_data_filtering_security_profile"
subcategory: "Objects"
---

# panos_data_filtering_security_profile

Manages data filtering security profiles.


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
resource "panos_data_filtering_security_profile" "example" {
    name = "example"
    description = "made by Terraform"
    rule {
        data_pattern = panos_custom_data_pattern_object.my_custom_obj.name
        applications = ["any"]
        file_types = ["any"]
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
* `data_capture` - (bool) Data capture.
* `rule` - (repeatable) Rule list spec, as defined below.

`rule` supports the following arguments:

* `name` - (Computed) Name.
* `data_pattern` - (Required) The data pattern name.
* `applications` - (list) List of applications.
* `file_types` - (list) List of file types.
* `direction` - Direction.  Valid values are `both` (default),
  `download`, or `upload`.
* `alert_threshold` - (int) Alert threshold.
* `block_threshold` - (int) Block threshold.
* `log_severity` - (PAN-OS 8.0+) Log severity.
