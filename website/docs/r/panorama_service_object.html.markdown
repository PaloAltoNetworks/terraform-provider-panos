---
layout: "panos"
page_title: "panos: panos_panorama_service_object"
sidebar_current: "docs-panos-panorama-resource-service-object"
description: |-
  Manages Panorama service objects.
---

# panos_panorama_service_object

This resource allows you to add/update/delete Panorama service objects.


## Import Name

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_service_object" "example" {
    name = "my_service"
    protocol = "tcp"
    description = "My service object"
    source_port = "2000-2049,2051-2099"
    destination_port = "32123"
    tags = ["internal", "dmz"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The service object's name.
* `device_group` - (Optional) The device group to put the service object into
  (default: `shared`).
* `description` - (Optional) The service object's description.
* `protocol` - (Required) The service's protocol.  This should be `tcp` or
  `udp`.
* `source_port` - (Optional) The source port.  This can be a single port
  number, range (1-65535), or comma separated (80,8080,443).
* `destination_port` - (Required) The destination port.  This can be a single
  port number, range (1-65535), or comma separated (80,8080,443).
* `tags` - (Optional) List of administrative tags.
