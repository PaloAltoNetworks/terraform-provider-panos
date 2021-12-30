---
page_title: "panos: panos_application_group"
subcategory: "Objects"
---

# panos_application_group

This resource allows you to add/update/delete application groups.


## PAN-OS

NGFW


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_application_group" "example" {
    name = "myApp"
    applications = [
        panos_application_object.a1.name,
        panos_application_object.a2.name,
    ]
}

resource "panos_application_object" "a1" {
    name = "app1"
    ...
}

resource "panos_application_object" "a2" {
    name = "app2"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `vsys` - (Optional) The group's vsys (default: `vsys1`).
* `name` - (Required) The group's name.
* `applications` - (Optional) List of applications in this group.
