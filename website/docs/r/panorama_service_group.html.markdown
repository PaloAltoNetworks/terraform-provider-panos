---
layout: "panos"
page_title: "panos: panos_panorama_service_group"
sidebar_current: "docs-panos-panorama-resource-service-group"
description: |-
  Manages Panorama service groups.
---

# panos_panorama_service_group

This resource allows you to add/update/delete Panorama service groups.


## Import Name

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_service_group" "example" {
    name = "static ntp grp"
    services = ["svc1", "svc2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service group's name.
* `device_group` - (Optional) The device group to put the service group into
  (default: `shared`).
* `services` - (Required) List of services to put in this service group.
* `tags` - (Optional) List of administrative tags.
