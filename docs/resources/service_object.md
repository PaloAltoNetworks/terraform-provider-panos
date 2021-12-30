---
page_title: "panos: panos_service_object"
subcategory: "Objects"
---

# panos_service_object

This resource allows you to add/update/delete service objects.


## PAN-OS

NGFW


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_service_object" "example" {
    name = "my_service"
    vsys = "vsys1"
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
* `vsys` - (Optional) The vsys to put the service object into (default:
  `vsys1`).
* `description` - (Optional) The service object's description.
* `protocol` - (Required) The service's protocol.  This should be `tcp`,
  `udp`, or `sctp` (PAN-OS 8.1+).
* `source_port` - (Optional) The source port.  This can be a single port
  number, range (1-65535), or comma separated (80,8080,443).
* `destination_port` - (Required) The destination port.  This can be a single
  port number, range (1-65535), or comma separated (80,8080,443).
* `tags` - (Optional) List of administrative tags.
* `override_session_timeout` - (Optional, bool, PAN-OS 8.1+) Set to `true` to
  override the default application timeouts.
* `override_timeout` - (Optional, int, PAN-OS 8.1+) The overridden TCP timeout.
* `override_half_closed_timeout` - (Optional, int, PAN-OS 8.1+) The overridden
  TCP half closed timeout.
* `override_time_wait_timeout` - (Optional, int, PAN-OS 8.1+) The overridden
  TCP wait time.
