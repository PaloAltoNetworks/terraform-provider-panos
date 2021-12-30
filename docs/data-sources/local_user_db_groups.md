---
page_title: "panos: panos_local_user_db_groups"
subcategory: "Device"
---

# panos_local_user_db_groups

Gets the list of local user database groups.


## Example Usage

```hcl
data "panos_local_user_db_groups" "example" {}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW:

* `vsys` - The vsys (default: `shared`).


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
