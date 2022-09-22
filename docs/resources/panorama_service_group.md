---
page_title: "panos: panos_panorama_service_group"
subcategory: "Objects"
---

# panos_panorama_service_group

This resource allows you to add/update/delete Panorama service groups.


## PAN-OS

Panorama


## Import Name

```shell
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_service_group" "example" {
    name = "static ntp grp"
    services = [
        panos_panorama_service_object.o1.name,
        panos_panorama_service_object.o2.name,
    ]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_service_object" "o1" {
    name = "svc1"
    ...

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_service_object" "o2" {
    name = "svc2"
    ...

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service group's name.
* `device_group` - (Optional) The device group to put the service group into
  (default: `shared`).
* `services` - (Required) List of services to put in this service group.
* `tags` - (Optional) List of administrative tags.
