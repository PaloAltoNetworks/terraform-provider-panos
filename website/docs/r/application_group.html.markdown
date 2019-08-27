---
layout: "panos"
page_title: "panos: panos_application_group"
description: |-
  Manages application groups.
---

# panos_application_group

This resource allows you to add/update/delete application groups.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_application_group" "example" {
    name = "myApp"
    applications = [
        "app1",
        "app2",
    ]
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The group's vsys (default: `vsys1`).
* `name` - (Required) The group's name.
* `applications` - (Optional) List of applications in this group.
