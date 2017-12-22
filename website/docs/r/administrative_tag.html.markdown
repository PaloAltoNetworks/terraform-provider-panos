---
layout: "panos"
page_title: "PANOS: panos_administrative_tag"
sidebar_current: "docs-panos-resource-administrative-tag"
description: |-
  Manages administrative tags.
---

# panos_administrative_tag

This resource allows you to add/update/delete administrative tags.

## Example Usage

```hcl
resource "panos_administrative_tag" "example" {
    name = "tag1"
    vsys = "vsys2"
    color = "color5"
    comment = "Internal resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The administrative tag's name.
* `vsys` - (Optional) The vsys to put the administrative tag into (default:
  `vsys1`).
* `color` - (Optional) The tag's color.  This should be either an empty string
  (no color) or a string such as `color1` or `color15`.  Note that for maximum
  portability, you should limit color usage to `color16`, which was available
  in PANOS 6.1.  PANOS 8.1's colors go up to `color42`.
* `comment` - (Optional) The administrative tag's description.
