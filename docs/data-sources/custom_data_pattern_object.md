---
page_title: "panos: panos_custom_data_pattern_object"
subcategory: "Objects"
---

# panos_custom_data_pattern_object

Gets information on a custom data pattern object.


## Example Usage

```hcl
# Predefined example.
data "panos_custom_data_pattern_object" "predef_example" {
    name = panos_custom_data_pattern_object.predef.name
}

resource "panos_custom_data_pattern_object" "predef" {
    name = "ex1"
    description = "made by Terraform"
    type = "predefined"
    predefined_pattern {
        name = "social-security-numbers"
        file_types = ["docx", "xlsx"]
    }

    lifecycle {
        create_before_destroy = true
    }
}


# Regex example.
data "panos_custom_data_pattern_object" "regex_example" {
    name = panos_custom_data_pattern_object.regex.name
}

resource "panos_custom_data_pattern_object" "regex" {
    name = "ex2
    description = "made by Terraform"
    type = "regex"
    regex {
        name = "blah"
        file_types = ["docx", "doc", "text/html"]
        regex = "shin megami tensei"
    }

    lifecycle {
        create_before_destroy = true
    }
}


# File property example.
data "panos_custom_data_pattern_object" "file_prop_example" {
    name = panos_custom_data_pattern_object.file_prop.name
}

resource "panos_custom_data_pattern_object" "file_prop" {
    name = "ex3"
    description = "made by Terraform"
    type = "file-properties"
    file_property {
        name = "blah2"
        file_type = data.panos_predefined_dlp_file_type.pdf_keywords.name
        file_property = data.panos_predefined_dlp_file_type.pdf_keywords.file_types.0.properties.0.name
        property_value = "foo"
    }

    lifecycle {
        create_before_destroy = true
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


## Attribute Reference

The following attributes are available:

* `description` - (Optional) The description.
* `type` - The type.
* `predefined_pattern` - (`type`=`predefined`, repeatable) List of
  predefined pattern specs, as definded below.
* `regex` - (`type`=`regex`) List of regex specs, as defined below.
* `file_property` - (`type`=`file-properties`) List of file properties specs,
  as defined below.

`predefined_pattern` supports the following arguments:

* `name` - The name.
* `file_types` - (list) List of file types.

`regex` supports the following arguments:

* `name` - Name.
* `file_types` - (list) List of file types.
* `regex` - The regex.

`file_property` supports the following arguments:

* `name` - Name.
* `file_type` - The file type.
* `file_property` - File property.
* `property_value` - Property value.
