---
page_title: "panos: panos_predefined_threat"
subcategory: "Predefined"
---


# panos_predefined_threat

Gets the predefined threats.


## Example Usage

```hcl
data "panos_predefined_threat" "example" {}
```


## Argument Reference

The following arguments are supported:

* `name` - A specific threat ID / name.
* `threat_regex` - A regex to apply to the threat name.
* `threat_name` - An exact match against the threat name.
* `threat_type` - The threat type.  Valid values are `phone-home` (default)
  or `vulnerability`.


## Attribute Reference

The following attributes are supported.

* `total` - (int) The number of file types.
* `threats` - List of matched threats specs, as defined below.

`threats` supports the following attributes:

* `name` - The threat name / ID.
* `threat_name` - The threat name.
