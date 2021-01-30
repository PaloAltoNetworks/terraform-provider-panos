---
page_title: "panos: panos_predefined_dlp_file_type"
subcategory: "Predefined"
---


# panos_predefined_dlp_file_type

Gets the predefined DLP file type information.


## Example Usage

```hcl
data "panos_predefined_dlp_file_type" "example" {}
```


## Argument Reference

The following arguments are supported:

* `name` - Filter on this specific name.
* `label` - A specific label to filter on.

## Attribute Reference

The following attributes are supported.

* `total` - (int) The number of file types.
* `file_types` - List of file types structs, as defined below.

`file_types` supports the following attributes:

* `name` - The file type.
* `properties` - List of property specs, as defined below.

`file_types.properties` supports the following attributes:

* `name` - The DLP property name.
* `label` - The DLP property label.
