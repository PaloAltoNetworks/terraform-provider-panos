---
page_title: "panos: panos_service_group"
subcategory: "Objects"
---

# panos_service_group

This resource allows you to add/update/delete service groups.


## PAN-OS

NGFW


## Import Name

```shell
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_service_group" "example" {
    name = "static ntp grp"
    services = [
        panos_service_object.o1.name,
    ]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_service_object" "o1" {
    name = "my_service"
    protocol = "tcp"
    source_port = "2000-2049"
    destination_port = "32123"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service group's name.
* `vsys` - (Optional) The vsys to put the service group into (default:
  `vsys1`).
* `services` - (Required) List of services to put in this service group.
* `tags` - (Optional) List of administrative tags.
