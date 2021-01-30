---
page_title: "panos: panos_predefined_tdb_file_type"
subcategory: "Predefined"
---


# panos_predefined_tdb_file_type

Gets the predefined TDB file type information.


## Example Usage

```hcl
data "panos_predefined_tdb_file_type" "example" {}
```


## Argument Reference

The following arguments are supported:

* `name` - Filter on this specific file type.
* `full_name` - The full name.
* `full_name_regex` - A regex to match against the full name.
* `data_ident_only` - (bool) Limit results to those with data_ident=`true`.


## Attribute Reference

The following attributes are supported.

* `total` - (int) The number of file types.
* `file_types` - List of file types structs, as defined below.

`file_types` supports the following attributes:

* `name` - The file type.
* `file_type_id` - (int) The ID
* `threat_name` - The threat name
* `full_name` - The full name.
* `data_ident` - (bool) Data ident
* `file_type_ident` - (bool) File type ident
