---
layout: "panos"
page_title: "panos: panos_panorama_application_group"
description: |-
  Manages Panorama application groups.
---

# panos_panorama_application_group

This resource allows you to add/update/delete Panorama application groups.


## Import Name

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_application_group" "example" {
    name = "myApp"
    applications = [
        "app1",
        "app2",
    ]
}
```

## Argument Reference

The following arguments are supported:

* `device_group` - (Optional) The group's device group (default: `shared`).
* `name` - (Required) The group's name.
* `applications` - (Optional) List of applications in this group.
