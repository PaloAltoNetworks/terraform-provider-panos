---
layout: "panos"
page_title: "panos: panos_service_group"
sidebar_current: "docs-panos-resource-service-group"
description: |-
  Manages service groups.
---

# panos_service_group

This resource allows you to add/update/delete service groups.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_service_group" "example" {
    name = "static ntp grp"
    services = ["svc1", "svc2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service group's name.
* `vsys` - (Optional) The vsys to put the service group into (default:
  `vsys1`).
* `services` - (Required) List of services to put in this service group.
* `tags` - (Optional) List of administrative tags.
