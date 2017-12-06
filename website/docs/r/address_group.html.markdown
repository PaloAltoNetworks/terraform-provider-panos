---
layout: "panos"
page_title: "PANOS: panos_address_group"
sidebar_current: "docs-panos-resource-address-group"
description: |-
  Manages address groups.
---

# panos_address_group

This resource allows you to add/update/delete address groups.

Address groups are either statically defined or dynamically defined, so only
`static` or `dynamic` should be defined within a given address group.

## Example Usage

```hcl
# Static group
resource "panos_address_group" "static1" {
    name = "static ntp grp"
    description = "My NTP servers"
    static = ["ntp1", "ntp2", "ntp3"]
}

# Dynamic group
resource "panos_address_group" "dag1" {
    name = "dynamic grp"
    description = "My internal NTP servers"
    dynamic = "'internal' and 'ntp'"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The address group's name.
* `vsys` - (Optional) The vsys to put the address group into (default:
  `vsys1`).
* `static` - (Optional) The address objects to include in this statically
  defined address group.
* `dynamic` - (Optional) The IP tags to include in this DAG.
* `description` - (Optional) The address group's description.
* `tag` - (Optional) List of administrative tags.
