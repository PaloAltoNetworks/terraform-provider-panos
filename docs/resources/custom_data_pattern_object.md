---
page_title: "panos: panos_custom_data_pattern_object"
subcategory: "Objects"
---

# panos_custom_data_pattern_object

Manages custom data pattern objects.


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
# Predefined example.
# For the names, these seem to be hardcoded per PAN-OS version and are not
# otherwise queriable.  It is recommended that you use
# `data.panos_custom_data_pattern_object` to determine what the underlying
# `name` param is for a given setting.
resource "panos_custom_data_pattern_object" "predef" {
    name = "ex1"
    description = "made by Terraform"
    type = "predefined"
    predefined_pattern {
        name = "social-security-numbers"
        file_types = ["docx", "xlsx"]
    }
}

# Regex example.
resource "panos_custom_data_pattern_object" "predef" {
    name = "ex2
    description = "made by Terraform"
    type = "regex"
    regex {
        name = "blah"
        file_types = ["docx", "doc", "text/html"]
        regex = "shin megami tensei"
    }
}

# File property example.
# Here again, you can either use `data.panos_custom_data_pattern_object` to see
# how something is currently configured or you can use the
# `data.panos_predefined_dlp_file_type` data source to discover a setting that's
# already configured.
resource "panos_custom_data_pattern_object" "predef" {
    name = "ex3"
    description = "made by Terraform"
    type = "file-properties"
    file_property {
        name = "blah2"
        file_type = data.panos_predefined_dlp_file_type.pdf_keywords.name
        file_property = data.panos_predefined_dlp_file_type.pdf_keywords.file_types.0.properties.0.name
        property_value = "foo"
    }
}

data "panos_predefined_dlp_file_type" "pdf_keywords" {
    name = "pdf"
    label = "Keywords"
}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.
* `description` - (Optional) The description.
* `type` - The type.  Valid values are `file-properties` (default),
  `predefined`, or `regex`.
* `predefined_pattern` - (`type`=`predefined`, repeatable) List of
  predefined pattern specs, as definded below.
* `regex` - (`type`=`regex`) List of regex specs, as defined below.
* `file_property` - (`type`=`file-properties`) List of file properties specs,
  as defined below.

`predefined_pattern` supports the following arguments:

* `name` - (Required) The name.
* `file_types` - (list) List of file types.

`regex` supports the following arguments:

* `name` - (Required) Name.
* `file_types` - (list) List of file types.
* `regex` - (Required) The regex.

`file_property` supports the following arguments:

* `name` - (Required) Name.
* `file_type` - (Required) The file type.
* `file_property` - (Required) File property.
* `property_value` - (Required) Property value.
