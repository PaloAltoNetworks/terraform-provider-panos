---
page_title: "panos: panos_panorama_application_group"
subcategory: "Objects"
---

# panos_panorama_application_group

This resource allows you to add/update/delete Panorama application groups.


## PAN-OS

Panorama


## Import Name

```shell
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_application_group" "example" {
    name = "myApp"
    applications = [
        panos_panorama_application_group.g1.name,
        panos_panorama_application_group.g2.name,
    ]
}

resource "panos_panorama_application_object" "g1" {
    name = "app1"
    ...
}

resource "panos_panorama_application_object" "g2" {
    name = "app2"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `device_group` - (Optional) The group's device group (default: `shared`).
* `name` - (Required) The group's name.
* `applications` - (Optional) List of applications in this group.
